package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/anonutopia/gowaves"
	"github.com/wavesplatform/gowaves/pkg/crypto"
)

type WavesMonitor struct {
	StartedTime int64
}

func (wm *WavesMonitor) start() {
	wm.StartedTime = time.Now().Unix() * 1000
	for {
		// todo - make sure that everything is ok with 100 here
		pages, err := gowaves.WNC.TransactionsAddressLimit(TokenAddress, 100)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}

		if len(pages) > 0 {
			for _, t := range pages[0] {
				wm.checkTransaction(&t)
			}
		}

		time.Sleep(time.Second * WavesMonitorTick)
	}
}

func (wm *WavesMonitor) checkTransaction(talr *gowaves.TransactionsAddressLimitResponse) {
	tr := &Transaction{TxID: talr.ID}
	db.FirstOrCreate(tr, tr)
	if !tr.Processed {
		wm.processTransaction(tr, talr)
	}
}

func (wm *WavesMonitor) processTransaction(tr *Transaction, talr *gowaves.TransactionsAddressLimitResponse) {
	attachment := ""

	if len(talr.Attachment) > 0 {
		attachment = string(crypto.MustBytesFromBase58(talr.Attachment))
	}

	if talr.Type == 7 {
		wm.processExchangeOrder(tr, talr)
	} else if talr.Type == 4 && talr.Recipient == TokenAddress && !exclude(conf.Exclude, talr.Sender) {
		if len(talr.AssetID) == 0 {
			wm.purchaseAsset(talr)
		} else if talr.AssetID == TokenID {
			wm.sellAsset(talr)
		} else if attachment == "collect" {
			wm.collectEarnings(talr)
		} else if talr.AssetID == AHRKId {
			wm.purchaseAssetAHRK(talr)
		}
	}

	tr.Processed = true
	db.Save(tr)
}

func (wm *WavesMonitor) purchaseAsset(talr *gowaves.TransactionsAddressLimitResponse) {
	messageTelegram(fmt.Sprintf("We just had a new AINT purchase: %.8f WAVES 🚀", float64(talr.Amount)/float64(SatInBTC)), TelAnonTeam)
	waves := talr.Amount - WavesFee - WavesExchangeFee
	if waves > 0 {
		a, p := wm.calculateAssetAmount(uint64(waves))
		abr, err := gowaves.WNC.AddressesBalance(TokenAddress)
		if err == nil {
			nabr, _ := gowaves.WNC.AddressesBalance(TokenAddress)
			if purchaseAsset(a, uint64(waves), TokenID, p) == nil {
				for abr.Balance == nabr.Balance {
					time.Sleep(time.Second * 10)
					nabr, _ = gowaves.WNC.AddressesBalance(TokenAddress)
				}

				sendAsset(a, TokenID, talr.Sender)
				wm.splitWaves(waves, talr.Sender)
			}
		}
	}
}

func (wm *WavesMonitor) purchaseAssetAHRK(talr *gowaves.TransactionsAddressLimitResponse) {
	messageTelegram(fmt.Sprintf("We just had a new AINT purchase: %.2f AHRK 🚀", float64(talr.Amount)/float64(AHRKDec)), TelAnonTeam)
	waves := talr.Amount * 100 / int(pc.Prices.HRK)
	a, _ := wm.calculateAssetAmount(uint64(waves))
	sendAsset(a, TokenID, talr.Sender)
}

func (wm *WavesMonitor) sellAsset(talr *gowaves.TransactionsAddressLimitResponse) {
	log.Printf("%#v\n\n", talr)
	logTelegram(fmt.Sprintf("%#v\n\n", talr))
}

func (wm *WavesMonitor) processExchangeOrder(tr *Transaction, talr *gowaves.TransactionsAddressLimitResponse) {
	if talr.Order1.Sender != talr.Order2.Sender {
		waves := int(float64(talr.Amount) / float64(SatInBTC) * float64(talr.Price))

		if talr.Order1.Sender != TokenAddress && talr.Order1.OrderType == "buy" {
			wm.splitWaves(waves, talr.Order1.Sender)
		}
	}
}

func (wm *WavesMonitor) splitWaves(waves int, sender string) {
	// 40% to founder
	founder := &User{Address: &conf.Founder}
	db.FirstOrCreate(founder, founder)
	founder.AmountWaves += uint(float64(waves) * 0.4)
	db.Save(founder)

	// 40% for buy offers
	kv := &KeyValue{Key: "buyFund"}
	db.FirstOrCreate(kv, kv)
	kv.ValueInt += uint64(float64(waves) * 0.4)
	db.Save(kv)

	// 5% to user who referred
	u := getUser(sender)
	r := &User{}
	db.First(r, u.ReferralID)
	r.AmountWaves += uint(float64(waves) * 0.05)
	db.Save(r)

	// 15% to AINTs holders (more than 1.0 AINT)
	rest := uint(float64(waves) * 0.15)
	ns, err := gowaves.WNC.NodeStatus()
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return
	}

	t, err := total(0, ns.BlockchainHeight-1, "")
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return
	}

	err = wm.doPayouts(ns.BlockchainHeight-1, "", t, int(rest))
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}
}

func (wm *WavesMonitor) doPayouts(height int, after string, total int, value int) error {
	abdr, err := gowaves.WNC.AssetsBalanceDistribution(TokenID, height, 100, after)
	if err != nil {
		return err
	}

	for a, v := range abdr.Items {
		if !exclude(conf.Exclude, a) && v > int(SatInBTC) {
			ratio := float64(v) / float64(total)
			amount := int(float64(value) * ratio)

			if amount > 0 {
				u := &User{Address: &a}
				if err := db.FirstOrCreate(u, u).Error; err != nil {
					tid := int(SatInBTC) + rand.Intn(999999999-int(SatInBTC))
					u.TelegramID = &tid
					db.FirstOrCreate(u, u)
				}
				u.AmountWaves += uint(amount)
				db.Save(u)
				log.Printf("Added interest: %s - %.8f", *u.Address, float64(amount)/float64(SatInBTC))
				logTelegram(fmt.Sprintf("Added interest: %s - %.8f", *u.Address, float64(amount)/float64(SatInBTC)))
			}
		}
	}

	if abdr.HasNext {
		return wm.doPayouts(height, abdr.LastItem, total, value)
	}

	return nil
}

func (wm *WavesMonitor) calculateAssetAmount(wavesAmount uint64) (amount uint64, price uint64) {
	opr, err := gowaves.WMC.OrderbookPair(TokenID, "WAVES", 10)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return 0, 0
	}

	waves := uint64(0)

	for i, a := range opr.Asks {
		if i == 0 {
			price = a.Price
		}
		if wavesAmount > 0 {
			w := a.Amount * a.Price / SatInBTC
			newWaves := uint64(0)
			if w < wavesAmount {
				newWaves = w
				amount += a.Amount
				waves += newWaves
				wavesAmount -= newWaves
			} else {
				newWaves = wavesAmount
				amount += uint64(float64(wavesAmount) / float64(a.Price) * float64(SatInBTC))
				waves += newWaves
				wavesAmount -= newWaves
			}
		}
	}

	return amount, price
}

func (wm *WavesMonitor) collectEarnings(talr *gowaves.TransactionsAddressLimitResponse) {
	u := getUser(talr.Sender)
	if u.ID != 0 && u.AmountWaves > 0 {
		err := sendAsset(uint64(u.AmountWaves), "", talr.Sender)
		if err == nil {
			u.AmountWaves = 0
			db.Save(u)
		}
	}
}

func initWavesMonitor() {
	wm = &WavesMonitor{}
	go wm.start()
}

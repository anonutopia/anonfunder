package main

import (
	"fmt"
	"log"
	"math"
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
		gowaves.WNC.Host = WavesNodeURL
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

	if talr.Timestamp >= wm.StartedTime {
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
			} else if talr.AssetID == USDNId {
				wm.purchaseAssetUSDN(talr)
			}
		}
	}

	tr.Processed = true
	db.Save(tr)
}

func (wm *WavesMonitor) purchaseAsset(talr *gowaves.TransactionsAddressLimitResponse) {
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
			}
		}
	}
}

func (wm *WavesMonitor) purchaseAssetAHRK(talr *gowaves.TransactionsAddressLimitResponse) {
	messageTelegram(fmt.Sprintf(tr("purchaseAhrk", "hr"), float64(talr.Amount)/float64(AHRKDec)), TelAnonTeam)
	waves := int(float64(talr.Amount) / float64(AHRKDec) / pc.getHRK() * float64(SatInBTC))
	a, _ := wm.calculateAssetAmount(uint64(waves))
	sendAsset(a, TokenID, talr.Sender)
}

func (wm *WavesMonitor) purchaseAssetUSDN(talr *gowaves.TransactionsAddressLimitResponse) {
	messageTelegram(fmt.Sprintf(tr("purchaseUsdn", "hr"), float64(talr.Amount)/float64(AHRKDec)), TelAnonTeam)
	waves := int(float64(talr.Amount) / float64(AHRKDec) / pc.Prices.USD * float64(SatInBTC))
	a, _ := wm.calculateAssetAmount(uint64(waves))
	sendAsset(a, TokenID, talr.Sender)
}

func (wm *WavesMonitor) sellAsset(talr *gowaves.TransactionsAddressLimitResponse) {
	log.Printf("%#v\n\n", talr)
	logTelegram(fmt.Sprintf("%#v\n\n", talr))
}

func (wm *WavesMonitor) processExchangeOrder(tra *Transaction, talr *gowaves.TransactionsAddressLimitResponse) {
	var priceChanged bool
	var newPrice float64
	var message string
	var messageEn string

	if talr.Order1.OrderType == "buy" &&
		talr.Order1.AssetPair.AmountAsset == TokenID {

		waves := int(((float64(talr.Amount) / float64(SatInBTC)) * (float64(talr.Price) / float64(SatInBTC))) * float64(SatInBTC))
		_, p := wm.calculateAssetAmount(uint64(waves))
		priceChanged, newPrice = wm.checkPriceRecord(p)

		amountEur := (float64(waves) / float64(SatInBTC)) * pc.Prices.EUR

		if talr.Order2.Sender == TokenAddress {
			wm.splitWaves(waves, talr.Order1.Sender)
			message = fmt.Sprintf(tr("purchase", "hr"), float64(waves)/float64(SatInBTC), amountEur)
			messageEn = fmt.Sprintf(tr("purchase", "en"), float64(waves)/float64(SatInBTC), amountEur)
		} else {
			messagePvt := fmt.Sprintf(tr("purchase", "hr"), float64(waves)/float64(SatInBTC), amountEur)
			messageTelegram(messagePvt, TelAnonTeam)
		}

		if priceChanged {
			if len(message) > 0 {
				message += "\n\n"
				messageEn += "\n\n"
			}
			message += fmt.Sprintf(tr("newAintPrice", "hr"), newPrice)
			messageEn += fmt.Sprintf(tr("newAintPrice", "en"), newPrice)
		}

		if len(message) > 0 {
			if priceChanged {
				messageTelegramPin(message, TelKriptokuna)
				messageTelegramPin(message, TelAnonBalkan)
				messageTelegramPin(messageEn, TelAnonutopia)
			} else {
				messageTelegram(message, TelKriptokuna)
				messageTelegram(message, TelAnonBalkan)
				messageTelegram(messageEn, TelAnonutopia)
			}
		}
	}
}

func (wm *WavesMonitor) splitWaves(waves int, sender string) {
	var rest uint
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
	if r.ID != 0 {
		db.First(r, u.ReferralID)
		r.AmountWaves += uint(float64(waves) * 0.05)
		db.Save(r)

		// 15% to AINTs holders (more than 1.0 AINT)
		rest = uint(float64(waves) * 0.15)
	} else {
		// 20% to AINTs holders (more than 1.0 AINT)
		rest = uint(float64(waves) * 0.2)
	}

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

			if amount > 0 && len(a) > 0 {
				u := getUser(a)
				if u.ID == 0 {
					u = um.createUserWeb(a)
				}
				newAmount := uint(float64(amount) / float64(SatInBTC) * pc.getHRK() * float64(AHRKDec))
				u.AmountAhrkAint += newAmount
				db.Save(u)
				log.Printf("Added interest: %s - %.6f", *u.Address, float64(newAmount)/float64(AHRKDec))
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

	for _, a := range opr.Asks {
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
			price = a.Price
		}
	}

	return amount, price
}

func (wm *WavesMonitor) collectEarnings(talr *gowaves.TransactionsAddressLimitResponse) {
	u := getUser(talr.Sender)
	if u.ID != 0 && u.AmountAhrkAint > 0 {
		err := sendAsset(uint64(u.AmountAhrkAint), AHRKId, talr.Sender)
		if err == nil {
			u.AmountAhrkAint = 0
			db.Save(u)
		}
	}
}

func (wm *WavesMonitor) checkPriceRecord(price uint64) (changed bool, newPrice float64) {
	kv := &KeyValue{Key: "aintPriceRecord"}
	db.FirstOrCreate(kv, kv)

	wPrice := float64(price) / float64(SatInBTC) * pc.Prices.EUR
	wPrice = math.Floor(wPrice*100) / 100
	wPriceInt := uint64(wPrice * 100)

	if kv.ValueInt < wPriceInt {
		kv.ValueInt = wPriceInt
		db.Save(kv)
		return true, wPrice
	}

	return false, wPrice
}

func initWavesMonitor() {
	wm = &WavesMonitor{}
	go wm.start()
}

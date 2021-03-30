package main

import (
	"log"
	"time"

	"github.com/anonutopia/gowaves"
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
	if talr.Type == 7 {
		wm.processExchangeOrder(tr, talr)
	} else if talr.Type == 4 && talr.Recipient == TokenAddress {
		if len(talr.AssetID) == 0 {
			wm.purchaseAsset(talr)
		} else if talr.AssetID == TokenID {
			wm.sellAsset(talr)
		}
	}
	tr.Processed = true
	db.Save(tr)
}

func (wm *WavesMonitor) purchaseAsset(talr *gowaves.TransactionsAddressLimitResponse) {
	log.Printf("%#v", talr)
}

func (wm *WavesMonitor) sellAsset(talr *gowaves.TransactionsAddressLimitResponse) {
	log.Printf("%#v", talr)
}

func (wm *WavesMonitor) processExchangeOrder(tr *Transaction, talr *gowaves.TransactionsAddressLimitResponse) {
	log.Printf("%#v", talr)
}

func initWavesMonitor() {
	wm = &WavesMonitor{}
	go wm.start()
}

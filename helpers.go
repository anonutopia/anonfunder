package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/anonutopia/gowaves"
	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

func sendAsset(amount uint64, assetId string, recipient string) error {
	if conf.Dev || conf.Debug {
		err := errors.New(fmt.Sprintf("Not sending (dev): %d - %s - %s", amount, assetId, recipient))
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	var assetBytes []byte

	// Create sender's public key from BASE58 string
	sender, err := crypto.NewPublicKeyFromBase58(conf.PublicKey)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Create sender's private key from BASE58 string
	sk, err := crypto.NewSecretKeyFromBase58(conf.PrivateKey)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Current time in milliseconds
	ts := time.Now().Unix() * 1000

	if len(assetId) > 0 {
		assetBytes = crypto.MustBytesFromBase58(assetId)
	} else {
		assetBytes = []byte{}
	}

	asset, err := proto.NewOptionalAssetFromBytes(assetBytes)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	assetW, err := proto.NewOptionalAssetFromBytes([]byte{})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	rec, err := proto.NewAddressFromString(recipient)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	tr := proto.NewUnsignedTransferWithSig(sender, *asset, *assetW, uint64(ts), amount, WavesFee, proto.Recipient{Address: &rec}, nil)

	err = tr.Sign(proto.MainNetScheme, sk)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Create new HTTP client to send the transaction to public TestNet nodes
	client, err := client.NewClient(client.Options{BaseUrl: WavesNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Context to cancel the request execution on timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// // Send the transaction to the network
	_, err = client.Transactions.Broadcast(ctx, tr)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	return nil
}

func purchaseAsset(amountAsset uint64, amountWaves uint64, assetId string, price uint64) error {
	if conf.Dev || conf.Debug {
		return errors.New(fmt.Sprintf("Not purchasing asset (dev): %d - %d - %s - %d", amountAsset, amountWaves, assetId, price))
	}

	var assetBytes []byte

	// Create sender's public key from BASE58 string
	sender, err := crypto.NewPublicKeyFromBase58(conf.PublicKey)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	matcher, err := crypto.NewPublicKeyFromBase58(MatcherPublicKey)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Create sender's private key from BASE58 string
	sk, err := crypto.NewSecretKeyFromBase58(conf.PrivateKey)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Current time in milliseconds
	ts := time.Now().Unix() * 1000
	ets := time.Now().Add(time.Hour*24*29).Unix() * 1000

	if len(assetId) > 0 {
		assetBytes = crypto.MustBytesFromBase58(assetId)
	} else {
		assetBytes = []byte{}
	}

	asset, err := proto.NewOptionalAssetFromBytes(assetBytes)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	assetW, err := proto.NewOptionalAssetFromBytes([]byte{})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	bo := proto.NewUnsignedOrderV1(sender, matcher, *asset, *assetW, proto.Buy, price, amountAsset, uint64(ts), uint64(ets), WavesExchangeFee)

	err = bo.Sign(proto.MainNetScheme, sk)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	_, err = gowaves.WMC.OrderbookMarketAlt(bo)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	return nil
}

func exclude(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func total(t int, height int, after string) (int, error) {
	abdr, err := gowaves.WNC.AssetsBalanceDistribution(TokenID, height, 100, after)
	if err != nil {
		return 0, err
	}

	for a, v := range abdr.Items {
		if !exclude(conf.Exclude, a) {
			t = t + v
		}
	}

	if abdr.HasNext {
		return total(t, height, abdr.LastItem)
	}

	return t, nil
}

type CalcResponse struct {
	Amount float64 `json:"amount"`
}

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Prices struct {
	BTC float64 `json:"BTC"`
	ETH float64 `json:"ETH"`
	EUR float64 `json:"EUR"`
	USD float64 `json:"USD"`
}

type PriceHNB struct {
	BrojTecajnice    string `json:"Broj tečajnice"`
	DatumPrimjene    string `json:"Datum primjene"`
	Drzava           string `json:"Država"`
	SifraValute      string `json:"Šifra valute"`
	Valuta           string `json:"Valuta"`
	Jedinica         int    `json:"Jedinica"`
	KupovniZaDevize  string `json:"Kupovni za devize"`
	SrednjiZaDevize  string `json:"Srednji za devize"`
	ProdajniZaDevize string `json:"Prodajni za devize"`
}

type PricesHNB struct {
	Prices []*PriceHNB
}

func (p *PricesHNB) getUSD() float64 {
	price := 0.0

	priceStr := strings.Replace(p.Prices[0].SrednjiZaDevize, ",", ".", -1)

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Println(err.Error())
		logTelegram(err.Error())
		return 0.0
	}

	return price
}

type PriceClient struct {
	Prices    *Prices
	PricesHNB *PricesHNB
}

func (pc *PriceClient) getHRK() float64 {
	return pc.Prices.USD * pc.PricesHNB.getUSD()
}

func (pc *PriceClient) doRequest() (*Prices, error) {
	p := &Prices{}
	cl := http.Client{}

	var req *http.Request
	var err error

	req, err = http.NewRequest(http.MethodGet, PricesURL, nil)

	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	res, err := cl.Do(req)

	if err == nil {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != 200 {
			log.Printf("[PriceClient.DoRequest] Error, body: %s", string(body))
			// logTelegram(fmt.Sprintf("[PriceClient.DoRequest] Error, body: %s", string(body)))
			err := errors.New(res.Status)
			return nil, err
		}
		json.Unmarshal(body, p)
	} else {
		return nil, err
	}

	return p, nil
}

func (pc *PriceClient) doRequestHNB() (*PricesHNB, error) {
	p := &PricesHNB{}
	cl := http.Client{}
	prices := make([]*PriceHNB, 0)

	var req *http.Request
	var err error

	req, err = http.NewRequest(http.MethodGet, PricesHNBURL, nil)

	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	res, err := cl.Do(req)

	if err == nil {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != 200 {
			log.Println(string(body))
			err := errors.New(res.Status)
			// logTelegram(err.Error())
			return nil, err
		}
		json.Unmarshal(body, &prices)
		p.Prices = prices
		pc.PricesHNB = p
	} else {
		return nil, err
	}

	return p, nil
}

func (pc *PriceClient) start() {
	go func() {
		for {
			if p, err := pc.doRequest(); err != nil {
				log.Printf("[PriceClient.start] Error: %s", err.Error())
				logTelegram(err.Error())
			} else {
				pc.Prices = p
			}

			if p, err := pc.doRequestHNB(); err != nil {
				log.Println(err.Error())
				// logTelegram(err.Error())
			} else {
				pc.PricesHNB = p
			}

			time.Sleep(time.Minute * 15)
		}
	}()
	pc.wait()
}

func (pc *PriceClient) wait() {
	for pc.Prices == nil {
		time.Sleep(time.Millisecond * 100)
	}
}

func initPriceClient() *PriceClient {
	pc := &PriceClient{}
	pc.start()
	return pc
}

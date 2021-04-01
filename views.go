package main

import (
	"fmt"
	"log"

	macaron "gopkg.in/macaron.v1"
)

func accumulatedEarnings(ctx *macaron.Context) string {
	ctx.Resp.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	address := ctx.Params("address")

	u := &User{Address: address}
	if err := db.Unscoped().FirstOrCreate(u, u).Error; err != nil {
		log.Println(err)
	}

	earnings := float64(u.AmountWaves) / float64(SatInBTC)
	response := fmt.Sprintf("document.getElementById('accumulatedEarnings').innerHTML = '%.4f';", earnings)
	return response
}

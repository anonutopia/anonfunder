package main

import (
	"fmt"

	macaron "gopkg.in/macaron.v1"
)

func accumulatedEarnings(ctx *macaron.Context) string {
	ctx.Resp.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	address := ctx.Params("address")

	u := &User{Address: address}
	db.FirstOrCreate(u, u)

	earnings := float64(u.AmountWaves) / float64(SatInBTC)
	response := fmt.Sprintf("document.getElementById('accumulatedEarnings').innerHTML = '%.8f';", earnings)
	return response
}

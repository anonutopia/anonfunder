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

	response := fmt.Sprintf(
		"document.getElementById('earningsWaves').value = %d;\n"+
			"document.getElementById('earningsAhrk').value = %d;\n"+
			"document.getElementById('earningsAeur').value = %d;\n",
		u.AmountWaves,
		u.AmountAhrk,
		u.AmountAeur,
	)

	return response
}

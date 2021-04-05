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

	if u.TempCode == nil || len(*u.TempCode) < 32 {
		rs := randString(32)
		u.TempCode = &rs
		db.Save(u)
	}

	res := "document.getElementById('earningsWaves').value = %d;\n" +
		"document.getElementById('earningsAhrk').value = %d;\n" +
		"document.getElementById('earningsAeur').value = %d;\n"

	if !u.FunderBotStarted {
		if address, err := encrypt([]byte(*u.TempCode), u.Address); err == nil {
			link := fmt.Sprintf("https://t.me/FunderBot?start=%s", address)
			res += "document.getElementById('btnFunderBot').classList.remove('disabled');\n"
			res += fmt.Sprintf("document.getElementById('btnFunderBot').href='%s';\n", link)
		} else {
			log.Println(err)
		}
	}

	response := fmt.Sprintf(
		res,
		u.AmountWaves,
		u.AmountAhrk,
		u.AmountAeur,
	)

	return response
}

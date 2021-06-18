package main

import (
	"fmt"
	"log"

	macaron "gopkg.in/macaron.v1"
)

func accumulatedEarnings(ctx *macaron.Context) string {
	ctx.Resp.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	address := ctx.Params("address")

	u := &User{Address: &address}
	if err := db.Unscoped().FirstOrCreate(u, u).Error; err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	ref := ctx.GetCookie("referral")
	if len(ref) > 0 {
		r := &User{}

		if err := db.Where("code = ?", ref).First(r).Error; err != nil {
			db.Where("nickname = ?", ref).First(r)
		}

		if r.ID != 0 && u.ReferralID == 0 {
			u.ReferralID = r.ID
		}
	}

	rs := randString(10)
	u.TempCode = &rs
	db.Save(u)

	res := "document.getElementById('earningsWaves').value = %d;\n" +
		"document.getElementById('earningsAhrk').value = %d;\n" +
		"document.getElementById('earningsAeur').value = %d;\n"

	if !u.FunderBotStarted {
		link := fmt.Sprintf("https://t.me/FunderRobot?start=%s", *u.TempCode)
		res += "document.getElementById('btnFunderBot').classList.remove('disabled');\n"
		res += fmt.Sprintf("document.getElementById('btnFunderBot').href='%s';\n", link)
	}

	// if !u.AnoteRobotStarted {
	// 	link := fmt.Sprintf("https://t.me/AnoteRobot?start=%s", *u.TempCode)
	// 	res += "document.getElementById('btnAnoteRobot').classList.remove('disabled');\n"
	// 	res += fmt.Sprintf("document.getElementById('btnAnoteRobot').href='%s';\n", link)
	// }

	response := fmt.Sprintf(
		res,
		u.AmountWaves,
		u.AmountAhrk,
		u.AmountAeur,
	)

	return response
}

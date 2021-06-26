package main

import (
	"fmt"
	"math"
	"strconv"

	macaron "gopkg.in/macaron.v1"
)

func accumulatedEarnings(ctx *macaron.Context) string {
	ctx.Resp.Header().Set("Content-Type", "application/javascript; charset=utf-8")

	address := ctx.Params("address")

	if len(address) == 0 {
		return ""
	}

	user := getUser(address)

	if user.ID == 0 {
		user = um.createUserWeb(address)
	}

	ref := ctx.GetCookie("referral")
	if len(ref) > 0 {
		r := &User{}

		if err := db.Where("code = ?", ref).First(r).Error; err != nil {
			db.Where("nickname = ?", ref).First(r)
		}

		if r.ID != 0 && user.ReferralID == 0 {
			user.ReferralID = r.ID
		}
	}

	rs := randString(10)
	user.TempCode = &rs
	db.Save(user)

	res := "document.getElementById('earningsWaves').value = %d;\n" +
		"document.getElementById('earningsAhrk').value = %d;\n" +
		"document.getElementById('earningsAeur').value = %d;\n" +
		"document.getElementById('referralLink').value += '%s';\n"

	if !user.FunderBotStarted {
		link := fmt.Sprintf("https://t.me/FunderRobot?start=%s", *user.TempCode)
		res += "document.getElementById('btnFunderBot').classList.remove('disabled');\n"
		res += fmt.Sprintf("document.getElementById('btnFunderBot').href='%s';\n", link)
	}

	// if !user.AnoteRobotStarted {
	// 	link := fmt.Sprintf("https://t.me/AnoteRobot?start=%s", *user.TempCode)
	// 	res += "document.getElementById('btnAnoteRobot').classList.remove('disabled');\n"
	// 	res += fmt.Sprintf("document.getElementById('btnAnoteRobot').href='%s';\n", link)
	// }

	response := fmt.Sprintf(
		res,
		user.AmountWaves,
		user.AmountAhrk,
		user.AmountAeur,
		*user.Code,
	)

	return response
}

func calculateAints(ctx *macaron.Context) {
	cr := &CalcResponse{}
	w := ctx.Params("waves")
	if wFloat, err := strconv.ParseFloat(w, 64); err == nil {
		wInt := uint64(wFloat * float64(SatInBTC))
		a, _ := wm.calculateAssetAmount(wInt)
		aFloat := float64(a) / float64(SatInBTC)
		amount := math.Floor(aFloat*float64(SatInBTC)) / float64(SatInBTC)
		cr.Amount = amount
	}

	ctx.JSON(200, cr)
}

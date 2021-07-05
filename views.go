package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/anonutopia/gowaves"
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

func websiteData(ctx *macaron.Context) string {
	ctx.Resp.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	response := ""
	aintPerc := 0.0
	circulatingAint := int64(0)
	ahrkPerc := 0.0
	circulatingAhrk := int64(0)

	res := "document.getElementById('limit-ahrk').style.width = '%.1f%%';\n" +
		"document.getElementById('circulating-ahrk').innerHTML = '%d';\n" +
		"document.getElementById('limit-ahrk2').innerHTML = '%.1f';\n" +
		"document.getElementById('limit-aint').style.width = '%.1f%%';\n" +
		"document.getElementById('circulating-aint').innerHTML = '%d';\n" +
		"document.getElementById('limit-aint1').innerHTML = '%.1f';\n" +
		"document.getElementById('limit-aint2').innerHTML = '%.1f';\n"

	abr, err := gowaves.WNC.AssetsBalance(AHRKAddress, AHRKId)
	if err == nil {
		circulatingAhrk = 1000000000000 - abr.Balance
		ahrkPerc = float64(circulatingAhrk) / float64(1000000000000)
		ahrkPerc = ahrkPerc * 100
	}

	abr, err = gowaves.WNC.AssetsBalance(TokenAddress, TokenID)
	if err == nil {
		circulatingAint = 1900000000000 - abr.Balance - 475000000000
		aintPerc = float64(circulatingAint) / float64(1900000000000)
		aintPerc = aintPerc * 100
	}

	response = fmt.Sprintf(
		res,
		ahrkPerc,
		circulatingAhrk/int64(AHRKDec),
		ahrkPerc,
		aintPerc,
		circulatingAint/int64(SatInBTC),
		aintPerc,
		aintPerc,
	)

	return response
}

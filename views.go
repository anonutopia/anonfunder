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
		user.AmountAhrkAint,
		user.AmountAhrk,
		user.AmountAeur,
		*user.Code,
	)

	return response
}

func calculateAints(ctx *macaron.Context) {
	wInt := uint64(0)
	cr := &CalcResponse{}
	c := ctx.Params("currency")
	a := ctx.Params("amount")
	if aFloat, err := strconv.ParseFloat(a, 64); err == nil {
		if c == "ahrk" {
			wInt = uint64(aFloat / pc.getHRK() * float64(SatInBTC))
		} else {
			wInt = uint64(aFloat * float64(SatInBTC))
		}
		aa, _ := wm.calculateAssetAmount(wInt)
		aFloat := float64(aa) / float64(SatInBTC)
		amount := math.Floor(aFloat*float64(SatInBTC)) / float64(SatInBTC)
		cr.Amount = amount
	}

	ctx.JSON(200, cr)
}

func websiteData(ctx *macaron.Context) string {
	ctx.Resp.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	response := ""
	// aintPerc := 0.0
	// circulatingAint := int64(0)
	aintcPerc := 0.0
	circulatingAintc := int64(0)
	ahrkPerc := 0.0
	circulatingAhrk := int64(0)

	res := "document.getElementById('limit-ahrk').style.width = '%.1f%%';\n" +
		"document.getElementById('circulating-ahrk').innerHTML = '%d';\n" +
		"document.getElementById('limit-ahrk2').innerHTML = '%.1f';\n" +
		"document.getElementById('limit-aintc').style.width = '%.1f%%';\n" +
		"document.getElementById('circulating-aintc').innerHTML = '%d';\n" +
		"document.getElementById('limit-aintc2').innerHTML = '%.1f';\n"

		// "document.getElementById('limit-aint').style.width = '%.1f%%';\n" +
		// "document.getElementById('circulating-aint').innerHTML = '%d';\n" +
		// "document.getElementById('limit-aint1').innerHTML = '%.1f';\n" +
		// "document.getElementById('limit-aint2').innerHTML = '%.1f';\n"

	abr, err := gowaves.WNC.AssetsBalance(AHRKAddress, AHRKId)
	if err == nil {
		abr2, err := gowaves.WNC.AssetsBalance(TokenAddress, AHRKId)
		if err == nil {
			abr3, err := gowaves.WNC.AssetsBalance("3PCGYBU7kG44GtXbZGUctCVcq9uR8W4eVXk", AHRKId)
			if err == nil {
				abr4, err := gowaves.WNC.AssetsBalance("3PLrCnhKyX5iFbGDxbqqMvea5VAqxMcinPW", AHRKId)
				if err == nil {
					circulatingAhrk = 1000000000000 - abr.Balance - abr2.Balance - abr3.Balance - abr4.Balance
					ahrkPerc = float64(circulatingAhrk) / float64(1000000000000)
					ahrkPerc = ahrkPerc * 100
				}
			}
		}
	}

	abr, err = gowaves.WNC.AssetsBalance(TokenAddress, TokenID)
	if err == nil {
		// circulatingAint = 1900000000000 - abr.Balance - 475000000000
		// aintPerc = float64(circulatingAint) / float64(1900000000000)
		// aintPerc = aintPerc * 100

		circulatingAintc = 1000000000000 - abr.Balance
		aintcPerc = float64(circulatingAintc) / float64(1000000000000)
		aintcPerc = aintcPerc * 100
	}

	response = fmt.Sprintf(
		res,
		ahrkPerc,
		circulatingAhrk/int64(AHRKDec),
		ahrkPerc,
		aintcPerc,
		circulatingAintc/int64(SatInBTC),
		aintcPerc,
		// aintPerc,
	)

	return response
}

package main

import (
	"github.com/bykovme/gotrans"
)

func initLangs() {
	gotrans.InitLocales("langs")
	gotrans.SetDefaultLocale("en")
}

func tr(text string, lang string) string {
	return gotrans.Tr(lang, text)
}

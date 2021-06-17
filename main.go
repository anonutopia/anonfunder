package main

import (
	"log"

	macaron "gopkg.in/macaron.v1"
	"gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

var conf *Config

var bot *telebot.Bot

var db *gorm.DB

var pc *PriceClient

var um *UserManager

var wm *WavesMonitor

var m *macaron.Macaron

func main() {
	initLangs()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	conf = initConfig()

	bot = initTelegramBot()

	db = initDb()

	pc = initPriceClient()

	um = initUserManager()

	initWavesMonitor()

	initCommands()

	m = initMacaron()

	m.Get("/:address/earnings.js", accumulatedEarnings)

	logTelegram("AnonFunder daemon successfully started. 🚀")
	log.Println("AnonFunder daemon successfully started. 🚀")

	bot.Start()
}

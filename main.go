package main

import (
	"log"

	"gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

var conf *Config

var bot *telebot.Bot

var db *gorm.DB

var pc *PriceClient

var um *UserManager

var wm *WavesMonitor

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

	logTelegram("AnonFunder daemon successfully started. 🚀")
	log.Println("AnonFunder daemon successfully started. 🚀")

	bot.Start()
}

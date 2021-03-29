package main

import (
	"log"
	"time"

	"gopkg.in/tucnak/telebot.v2"
	tb "gopkg.in/tucnak/telebot.v2"
)

func initTelegramBot() *telebot.Bot {
	b, err := tb.NewBot(tb.Settings{
		Token:     conf.TelegramAPIKey,
		Poller:    &tb.LongPoller{Timeout: TelPollerTimeout * time.Second},
		Verbose:   conf.Debug,
		ParseMode: tb.ModeHTML,
	})

	if err != nil {
		log.Fatal(err)
	}

	return b
}

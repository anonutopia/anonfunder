package main

import (
	"fmt"
	"log"

	"github.com/bykovme/gotrans"
	tb "gopkg.in/tucnak/telebot.v2"
)

func initCommands() {
	bot.Handle("/start", startCommand)
	bot.Handle("/status", statusCommand)
	bot.Handle(tb.OnText, unknownCommand)
}

func startCommand(m *tb.Message) {
	if len(m.Payload) > 0 {
		u := &User{TempCode: &m.Payload}
		if err := db.FirstOrCreate(u, u).Error; err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}

		if u.TelegramID == nil || *u.TelegramID == 0 {
			u.TelegramID = &m.Sender.ID
			u.FunderBotStarted = true
			if err := db.Save(u).Error; err != nil {
				log.Println(err)
				logTelegram(err.Error())
			}
		} else {
			u.FunderBotStarted = true
			if err := db.Save(u).Error; err != nil {
				log.Println(err)
				logTelegram(err.Error())
			}
		}
	} else {
		um.createUser(m)
	}

	bot.Send(m.Sender, gotrans.T("welcome"))
}

func statusCommand(m *tb.Message) {
	u := getUserByTelegramID(m)
	status := fmt.Sprintf(
		gotrans.T("status"),
		u.getAddress(),
		u.getAmountAint(),
		u.getAmountWaves(),
	)
	bot.Send(m.Sender, status)
}

func unknownCommand(m *tb.Message) {
	if m.Private() {
		bot.Send(m.Sender, gotrans.T("unknown"))
	}
}

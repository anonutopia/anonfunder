package main

import (
	"fmt"

	"github.com/bykovme/gotrans"
	tb "gopkg.in/tucnak/telebot.v2"
)

func initCommands() {
	bot.Handle("/start", startCommand)
	bot.Handle("/status", statusCommand)
	bot.Handle(tb.OnText, unknownCommand)
}

func startCommand(m *tb.Message) {
	um.createUser(m)
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
	bot.Send(m.Sender, gotrans.T("unknown"))
}

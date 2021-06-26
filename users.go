package main

import (
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

type UserManager struct {
}

func (um *UserManager) createUser(m *tb.Message) {
	u := &User{TelegramID: &m.Sender.ID}
	r := &User{}

	db.FirstOrCreate(u, u)
	db.First(r, m.Payload)

	if r.ID != 0 && r.ID != u.ID {
		u.Referral = r
		db.Save(u)
	}
}

func (um *UserManager) createUserWeb(address string) *User {
	u := &User{Address: &address}

	if err := db.FirstOrCreate(u, u).Error; err != nil {
		log.Println(err)
		logTelegram(err.Error())
	} else {
		rs := randString(10)
		u.Code = &rs
		if err := db.Save(u).Error; err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}
	}

	return u
}

func initUserManager() *UserManager {
	um := &UserManager{}
	return um
}

package main

import (
	"github.com/anonutopia/gowaves"
	"gorm.io/gorm"
)

// KeyValue model is used for storing key/values
type KeyValue struct {
	gorm.Model
	Key      string `gorm:"size:255;uniqueIndex"`
	ValueInt uint64
	ValueStr string
}

// User represents Telegram user
type User struct {
	gorm.Model
	Address     string `gorm:"size:255;uniqueIndex"`
	TelegramID  int    `gorm:"uniqueIndex"`
	ReferralID  uint
	Referral    *User
	AmountWaves uint
}

func getUser(address string) *User {
	u := &User{Address: address}
	db.First(u, u)
	return u
}

func (u *User) getAddress() string {
	if len(u.Address) > 0 {
		return u.Address
	}

	return "no address"
}

func (u *User) getAmountAint() float64 {
	abr, err := gowaves.WNC.AssetsBalance(u.Address, TokenID)
	if err == nil {
		return float64(abr.Balance) / float64(SatInBTC)
	}
	return 0
}

func (u *User) getAmountWaves() float64 {
	return float64(u.AmountWaves) / float64(SatInBTC)
}

// Transaction represents node's transaction
type Transaction struct {
	gorm.Model
	TxID      string `gorm:"size:255;uniqueIndex"`
	Processed bool   `gorm:"DEFAULT:false"`
}

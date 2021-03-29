package main

import (
	"gorm.io/gorm"
)

// KeyValue model is used for storing key/values
type KeyValue struct {
	gorm.Model
	Key      string `sql:"size:255;unique_index"`
	ValueInt uint64
	ValueStr string
}

// User represents Telegram user
type User struct {
	gorm.Model
	Address     string `sql:"size:255;unique_index"`
	TelegramID  int    `sql:"unique_index"`
	ReferralID  uint
	Referral    *User
	AmountAint  uint
	AmountWaves uint
	AmountBtc   uint
	AmountEth   uint
}

func (u *User) getAddress() string {
	if len(u.Address) > 0 {
		return u.Address
	}

	return "no address"
}

func (u *User) getAmountAint() float64 {
	return float64(u.AmountAint) / float64(SatInBTC)
}

func (u *User) getAmountWaves() float64 {
	return float64(u.AmountWaves) / float64(SatInBTC)
}

func (u *User) getAmountBtc() float64 {
	return float64(u.AmountBtc) / float64(SatInBTC)
}

func (u *User) getAmountEth() float64 {
	return float64(u.AmountEth) / float64(SatInBTC)
}

// Transaction represents node's transaction
type Transaction struct {
	gorm.Model
	TxID      string `sql:"size:255"`
	Processed bool   `sql:"DEFAULT:false"`
}

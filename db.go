package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initDb() *gorm.DB {
	var db *gorm.DB
	var err error
	dbconf := gorm.Config{}

	if conf.Debug {
		dbconf.Logger = logger.Default.LogMode(logger.Info)
	} else {
		dbconf.Logger = logger.Default.LogMode(logger.Error)
	}

	if conf.Dev {
		db, err = gorm.Open(sqlite.Open("anonfunder.db"), &dbconf)
	} else {
		db, err = gorm.Open(postgres.Open(conf.PostgreSQL), &dbconf)
	}

	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	if err := db.AutoMigrate(&KeyValue{}, &User{}, &Transaction{}); err != nil {
		panic(err.Error())
	}

	return db
}

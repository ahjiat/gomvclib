package Database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB = gorm.DB

var container = map[string]*gorm.DB{}

func AddDBConnection(sessionName string, dsn string) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}); if err != nil { panic(err) }
	/*
	sqlDB, err := db.DB(); if err != nil { panic(err) }
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(10)
	go func() {
		for i:=0;i<1000;i++  {
			stats := sqlDB.Stats()
			fmt.Println(stats)
			time.Sleep(time.Second)
		}
	}()
	*/
	container[sessionName] = db
}
func GetSession(sessionName string) *gorm.DB {
	return container[sessionName]
}

package Database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var container = map[string]*gorm.DB{}

func AddDBConnection(sessionName string, dsn string) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}); if err != nil { panic(err) }
	container[sessionName] = db
}
func GetSession(sessionName string) *gorm.DB {
	return container[sessionName]
}

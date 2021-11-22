package Database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var container = map[string]*sql.DB{}

func AddDBConnection(sessionName string, dns string) {
	db, err := sql.Open("mysql", dns); if err != nil { panic(err) }
	container[sessionName] = db
}
func GetSession(sessionName string) *sql.DB {
	return container[sessionName]
}

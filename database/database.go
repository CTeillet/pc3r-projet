package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	db, err := sql.Open("mysql", "projet-pc3r:QXKn6BmR7B7Iwfrs@tcp(cteillet.ddns.net:3306)/projet-pc3r?parseTime=true")

	if err != nil {
		panic(err.Error())
		return nil
	} else {
		return db
	}
}
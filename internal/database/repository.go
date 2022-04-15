package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

func PostTemperature(temp int) {
	db, err := sqlx.Connect("mysql", "admin:qwe123@tcp(localhost:3306)/test")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = db.Exec(
		"INSERT INTO temperatures (date, temp) VALUES(?,?)",
		time.Now().Add(time.Duration(time.Hour*7)),
		temp,
	)
	if err != nil {
		log.Println(err)
		return
	}
}

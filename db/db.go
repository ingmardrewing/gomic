package db

import (
	"database/sql"
	"log"

	"github.com/ingmardrewing/gomic/config"
)

var db *sql.DB

func Init() {
	dsn := config.GetDsn()
	log.Println(dsn)
	d, _ := sql.Open("mysql", dsn)
	err := d.Ping()
	if nil != err {
		panic(err)
	}
	db = d
}

func Query(query string) *sql.Rows {
	rows, err := db.Query(query)
	//defer rows.Close()
	if err != nil {
		panic(err.Error())
	}
	return rows
}

package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ingmardrewing/gomic/config"
	"github.com/ingmardrewing/gomic/page"
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

func InsertPage(p *page.Page) {
	ins := fmt.Sprintf("INSERT INTO pages VALUES('%s', '%s', '%s', '%s');\n", p.Title(), p.FSPath(), p.ImgUrl(), p.DisqusId())
	_, err := db.Exec(ins)
	if err != nil {
		panic(err.Error())
	}
}

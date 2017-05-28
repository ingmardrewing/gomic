package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ingmardrewing/gomic/comic"
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

func InsertPage(p *comic.Page) {
	ins := fmt.Sprintf("INSERT INTO pages (title, path, imgUrl, disqusId, act) VALUES('%s', '%s', '%s', '%s', '%s');\n", p.GetTitle(), p.FSPath(), p.GetImgUrl(), p.GetDisqusId(), "Act III")
	_, err := db.Exec(ins)
	if err != nil {
		panic(err.Error())
	}
}

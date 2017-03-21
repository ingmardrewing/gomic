package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ingmardrewing/gomic/comic"
	"github.com/ingmardrewing/gomic/config"
	"github.com/ingmardrewing/gomic/output"
)

func main() {
	config.Read("/Users/drewing/Sites/gomic.yaml")
	comic := comic.NewComic()
	output := output.NewOutput(&comic)
	output.WriteToFilesystem()

	user := os.Getenv("DB_GOMIC_USER")
	pass := os.Getenv("DB_GOMIC_PASS")
	name := os.Getenv("DB_GOMIC_NAME")
	host := os.Getenv("DB_GOMIC_HOST")
	dsn := fmt.Sprintf("%s:%s@%s/%s", user, pass, host, name)

	db, _ := sql.Open("mysql", dsn)

	db.Ping()
	Read(config.PngDir())
}

func Read(path string) {
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		fmt.Println(f.Name())
	}
}

package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ingmardrewing/gomic/config"
)

func main() {
	config.Read("/Users/drewing/Sites/gomic.yaml")
	//comic := comic.NewComic()
	//	output := output.NewOutput(&comic)
	//	output.WriteToFilesystem()

	db := config.ConnectDb()
	rows, err := db.Query("SELECT * FROM pages;")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var (
			title    sql.NullString
			path     sql.NullString
			imgUrl   sql.NullString
			disqusId sql.NullString
		)
		rows.Scan(&title, &path, &imgUrl, &disqusId)
		fmt.Printf("%s - %s - %s - %s \n", title.String, path.String, imgUrl.String, disqusId.String)
	}
	//Read(config.PngDir())

	//return &Page{title, path, imgUrl, servedrootpath, disqusId,
	/*
		for _, p := range config.Pages() {
			ins := fmt.Sprintf("INSERT INTO pages VALUES('%s', '%s', '%s', '%s');\n", p["title"], p["path"], p["imgUrl"], p["disqusId"])
			_, err := db.Exec(ins)
			if err != nil {
				panic(err.Error())
			}
		}
	*/
}

func Read(path string) {
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		fmt.Println(f.Name())
	}
}

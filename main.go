package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/ingmardrewing/gomic/comic"
	"github.com/ingmardrewing/gomic/config"
	"github.com/ingmardrewing/gomic/db"
	"github.com/ingmardrewing/gomic/fs"
)

func main() {
	config.Read("/Users/drewing/Sites/gomic.yaml")

	db.Init()
	rows := db.Query("SELECT * FROM pages;")

	comic := comic.NewComic(rows)

	currentImageFiles := fs.ReadImageFilenames()
	comic.CheckForNewPages(currentImageFiles)
	comic.ConnectPages()

	output := fs.NewOutput(&comic)
	output.WriteToFilesystem()

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

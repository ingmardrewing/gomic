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
}

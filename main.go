package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ingmardrewing/gomic/aws"
	"github.com/ingmardrewing/gomic/comic"
	"github.com/ingmardrewing/gomic/config"
	"github.com/ingmardrewing/gomic/db"
	"github.com/ingmardrewing/gomic/fs"
	"github.com/ingmardrewing/gomic/socmed"
	"github.com/ingmardrewing/gomic/strato"
)

func main() {
	config.Read("/Users/drewing/Sites/gomic.yaml")

	db.Init()
	rows := db.Query("SELECT * FROM pages order by pageNumber;")

	comic := comic.NewComic(rows)

	currentImageFiles := fs.ReadImageFilenames()
	checkForNewPages(currentImageFiles, comic)
	comic.ConnectPages()

	output := fs.NewOutput(&comic)
	output.WriteToFilesystem()

	if config.IsTest() {
		strato.UploadTest()
	} else if config.IsProd() {
		strato.UploadProd()
		socmed.Publish()
	}
}

func checkForNewPages(filenames []string, c comic.Comic) {
	for _, f := range filenames {
		if c.IsNewFile(f) {
			log.Printf("Found new file: %s", f)
			c.CreateThumbnail(f)
			p := comic.NewPageFromFilename(f)
			aws.UploadPage(p)
			db.InsertPage(p)
			c.AddPage(p)
			socmed.Prepare(p.Path(), p.Title(), p.ImgUrl(), p.ProdUrl(), p.Description())
		}
	}
}

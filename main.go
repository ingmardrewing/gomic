package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ingmardrewing/gomic/aws"
	"github.com/ingmardrewing/gomic/comic"
	"github.com/ingmardrewing/gomic/config"
	"github.com/ingmardrewing/gomic/db"
	"github.com/ingmardrewing/gomic/socmed"
	"github.com/ingmardrewing/gomicPages/content"
)

func main() {

	config.Read("/Users/drewing/Sites/gomic.yaml")

	db.Init()

	callApi()
	/*
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
			socmed.Publish(&comic)
		}
	*/
}

func callApi() {
	url := "https://drewing.eu/0.1/gomic/socmed/echo"
	data := []byte(`{"Link":"https://devabo.de/2017/05/27/87-Incoming","ImgUrl":"https://s3-us-west-1.amazonaws.com/devabode-us/comicstrips/DevAbode_0087.png","Title":"#87 Incoming","TagsCsvString":"comic,webcomic,graphicnovel,drawing,art,narrative,scifi,sci-fi,science-fiction,dystopy,parody,humor,nerd,pulp,geek,blackandwhite","Description":"While Eezer and Master Branch are talking about the losses from the attack ofthe cult, another problem occurs ..."}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("NewReqeust: ", err)
		return
	}

	user, pass := config.GetBasicAuthUserAndPass()
	req.SetBasicAuth(user, pass)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}
	defer resp.Body.Close()

	var t content.Page
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Println(t.Title)
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
			log.Printf("new File with Path: %s\n", p.Path())
			socmed.Prepare(p.Path(), p.Title(), p.ImgUrl(), p.ProdUrl(), p.Description())
		}
	}
}

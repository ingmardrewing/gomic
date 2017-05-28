package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ingmardrewing/gomic/aws"
	"github.com/ingmardrewing/gomic/comic"
	"github.com/ingmardrewing/gomic/config"
	"github.com/ingmardrewing/gomic/db"
	"github.com/ingmardrewing/gomic/socmed"
)

func main() {
	config.Read("/Users/drewing/Sites/gomic.yaml")
	db.Init()

	pages := callPagesApi()
	fmt.Println(pages)
	//	rows := db.Query("SELECT * FROM pages order by pageNumber;")

	//	comic := comic.NewComic(pages)

	/*
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

/*
[
  {
     "Id": 1,
    "PageNumber": 1,
   "Title": "#1 A Step in the dark",
      "Description": "#1 A Step in the dark",
     "Path": "/2013/08/01/a-step-in-the-dark",
    "ImgUrl": "https://devabode-us.s3.amazonaws.com/comicstrips/DevAbode_0001.png",
   "DisqusId": "8 http://devabo.de/?p=8",
      "Act": "Act I"
    },
	]
*/

type Pages struct {
	Pages []comic.Page
}

func callPagesApi() *Pages {
	url := "https://drewing.eu:8443/0.1/gomic/page/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	user, pass := config.GetBasicAuthUserAndPass()
	req.SetBasicAuth(user, pass)

	myClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := myClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	p := new(Pages)
	err = json.NewDecoder(resp.Body).Decode(p)
	if err != nil {
		panic(err)
	}
	return p
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
			log.Printf("new File with Path: %s\n", p.GetPath())
			socmed.Prepare(p.GetPath(), p.GetTitle(), p.GetImgUrl(), p.GetProdUrl(), p.GetDescription())
		}
	}
}

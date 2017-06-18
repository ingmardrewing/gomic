package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

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

	pages := callPagesApi()
	log.Println(pages)
	imageFiles := fs.ReadImageFilenames()
	newPages := checkForNewPages(imageFiles, pages.Pages)
	storeNewPages(newPages)
	pages.Pages = append(pages.Pages, newPages...)
	fmt.Println(newPages)

	comic := comic.NewComic(pages.Pages)
	comic.ConnectPages()

	output := fs.NewOutput(&comic)
	output.WriteToFilesystem()

	if config.IsTest() {
		strato.UploadTest()
	} else if config.IsProd() {
		strato.UploadProd()
		if len(newPages) > 0 {
			socmed.Publish(&comic)
		}
	}
}

func callPagesApi() *comic.Pages {
	url := "https://drewing.eu:8443/0.1/gomic/page/"
	log.Printf("callPagesApi with %s\n", url)
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

	p := new(comic.Pages)
	err = json.NewDecoder(resp.Body).Decode(p)
	if err != nil {
		panic(err)
	}
	return p
}

func checkForNewPages(filenames []string, knownPages []*comic.Page) []*comic.Page {
	newPages := []*comic.Page{}
	for _, f := range filenames {
		if isNewFile(f, knownPages) {
			comic.CreateThumbnail(f)
			p := comic.NewPageFromFilename(f)
			aws.UploadPage(p)
			db.InsertPage(p)
			newPages = append(newPages, p)
			log.Printf("new File with Path: %s\n", p.GetPath())
			socmed.Prepare(p.GetPath(), p.GetTitle(), p.GetImgUrl(), p.GetProdUrl(), p.GetDescription())
		}
	}
	return newPages
}

func isNewFile(filename string, knownPages []*comic.Page) bool {
	if !isRelevant(filename) {
		return false
	}
	for _, p := range knownPages {
		fn := p.GetImageFilename()
		if fn == filename {
			return false
		}
	}
	log.Println("File is new:" + filename)
	return true
}

func isRelevant(filename string) bool {
	irr := ".DS_Store"
	if filename == irr {
		return false
	}
	thumb := regexp.MustCompile(`^thumb_`)
	if thumb.MatchString(filename) {
		return false
	}
	return true
}

func storeNewPages(newPages []*comic.Page) {
	myClient := &http.Client{Timeout: 10 * time.Second}
	for _, p := range newPages {
		data := []byte(fmt.Sprintf(`{"Id":"%d","PageNumber":"%d","Title":"%s","Description":"%s","Path":"%s","ImgUrl":"%s","DisqusId":"%s","Act":"%s"}`,
			p.Id,
			p.PageNumber,
			p.Title,
			p.Description,
			p.Path,
			p.ImgUrl,
			p.DisqusId,
			p.Act))
		url := "https://drewing.eu:8443/0.1/gomic/page/"
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")

		if err != nil {
			panic(err)
		}

		user, pass := config.GetBasicAuthUserAndPass()
		req.SetBasicAuth(user, pass)

		resp, err := myClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("Put response Status:", resp.Status)
		fmt.Println("Put response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Put response Body:", body)
	}
}

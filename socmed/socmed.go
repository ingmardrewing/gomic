package socmed

import (
	"fmt"
	"log"

	"github.com/go-resty/resty"
)

var imgurl = ""
var prodUrl = ""
var title = ""
var path = ""

func Prepare(p string, t string, i string, pu string) {
	title = t
	path = p
	imgurl = i
	prodUrl = pu
}

func Publish() {
	response, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(getPublishableConted()).
		Post("http://user:}vh9m*,4#D4W7z2o2Hk%@drewing.eu:8081/0.1/gomic/socmed/publish")
	if err != nil {
		log.Println(err)
	}
	log.Println(response)
}

func getPublishableConted() string {
	tags := "comic,webcomic,graphicnovel,drawing,art,narrative,scifi,sci-fi,science-fiction,dystopy,parody,humor,nerd,pulp,geek,blackandwhite"
	description := "A new page of DevAbo.de is online ..."
	return fmt.Sprintf(`{"Link":"%s","ImgUrl":"%s","Title":"%s","TagsCsvString":"%s","Description":"%s"}`, prodUrl, imgurl, title, tags, description)
}

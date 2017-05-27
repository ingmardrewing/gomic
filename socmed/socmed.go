package socmed

import (
	"fmt"
	"log"

	"github.com/go-resty/resty"
	"github.com/ingmardrewing/gomic/config"
)

var imgurl = ""
var prodUrl = ""
var title = ""
var path = ""
var description = ""

func Prepare(p string, t string, i string, pu string, d string) {
	title = t
	path = p
	imgurl = i
	prodUrl = pu
	description = d
}

func Publish() {
	user, pass := config.GetBasicAuthUserAndPass()
	content := getPublishableConted()
	log.Println(content)
	response, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(getPublishableConted()).
		Post("https://" + user + ":" + pass + "@drewing.eu/0.1/gomic/socmed/publish")
	if err != nil {
		log.Println(err)
	}
	log.Println(response)
}

func getPublishableConted() string {
	tags := "comic,webcomic,graphicnovel,drawing,art,narrative,scifi,sci-fi,science-fiction,dystopy,parody,humor,nerd,pulp,geek,blackandwhite"
	return fmt.Sprintf(`{"Link":"%s","ImgUrl":"%s","Title":"%s","TagsCsvString":"%s","Description":"%s"}`, prodUrl, imgurl, title, tags, description)
}

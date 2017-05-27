package socmed

import (
	"fmt"
	"log"

	"github.com/ingmardrewing/gomic/comic"
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

func Publish(c *comic.Comic) {

	if notPrepared() {
		prepareFromComic(c)
	}
	user, pass := config.GetBasicAuthUserAndPass()
	content := getPublishableConted()

	log.Printf(`curl -X POST -H "Content-Type: application/json; charset=utf-8" -d '%s' -u %s:'%s' https://drewing.eu:443/0.1/gomic/socmed/publish`, content, user, pass)

	/*
		response, err := resty.R().
			SetHeader("Content-Type", "application/json").
			SetBody(getPublishableConted()).
			Post("https://" + user + ":" + pass + "@drewing.eu/0.1/gomic/socmed/publish")
		if err != nil {
			log.Println(err)
		}
		log.Println(response)
	*/
}

func prepareFromComic(c *comic.Comic) {
	lastPage := c.Get10LastComicPagesNewestFirst()[10]
	title = lastPage.Title()
	path = lastPage.Path()
	imgurl = lastPage.ImgUrl()
	prodUrl = lastPage.ProdUrl()
	description = lastPage.Description()
}

func notPrepared() bool {
	return len(prodUrl) == 0
}

func getPublishableConted() string {
	tags := "comic,webcomic,graphicnovel,drawing,art,narrative,scifi,sciencefiction,dystopy,parody,humor,nerd,pulp,geek,blackandwhite,inked,inking,illustration,drawings,draw,illustrator,sketch,sketchbook,artist,comics,comicart,concept,conceptart,create,creative,imagination,artwork,ink,fantasy"
	return fmt.Sprintf(`{"Link":"%s","ImgUrl":"%s","Title":"%s","TagsCsvString":"%s","Description":"%s"}`, prodUrl, imgurl, title, tags, description)
}

package socmed

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

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
	content := []byte(getPublishableConted())

	url := "https://drewing.eu:443/0.1/gomic/socmed/all/publish"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(content))
	req.SetBasicAuth(user, pass)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response status: ", resp.Status)
	fmt.Println("response headers: ", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response body: ", body)

	//log.Printf(`curl -X POST -H "Content-Type: application/json; charset=utf-8" -d '%s' -u %s:'%s' https://drewing.eu:443/0.1/gomic/socmed/all/publish`, content, user, pass)
}

func prepareFromComic(c *comic.Comic) {
	lastPage := c.Get10LastComicPagesNewestFirst()[10]
	title = lastPage.GetTitle()
	path = lastPage.GetPath()
	imgurl = lastPage.GetImgUrl()
	prodUrl = lastPage.GetProdUrl()
	description = lastPage.GetDescription()
}

func notPrepared() bool {
	return len(prodUrl) == 0
}

func getPublishableConted() string {
	tags := "webcomic,graphicnovel,comic,comicart,comics,sciencefiction,scifi,geek,nerd,art,artist,artwork,blackandwhite,concept,conceptart,create,creative,design,digital,draw,drawing,drawings,dystopy,fantasy,humor,illustration,illustrator,image,imagination,ink,inked,inking,kunst,malen,malerei,narrative,parody,pulp,sketch,sketchbook,tusche,zeichnen,zeichnung"
	return fmt.Sprintf(`{"Link":"%s","ImgUrl":"%s","Title":"%s","TagsCsvString":"%s","Description":"%s"}`, prodUrl, imgurl, title, tags, description)
}

package socmed

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ingmardrewing/gomic/comic"
	"github.com/ingmardrewing/gomic/config"
)

func getRequest(jsonData string) *http.Request {
	url := "https://drewing.eu:443/0.1/gomic/socmed/all/publish"
	user, pass := config.GetBasicAuthUserAndPass()
	jsonAsBytes := []byte(jsonData)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonAsBytes))
	req.SetBasicAuth(user, pass)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	return req
}

func Publish(c *comic.Comic) {
	jsonData := getJsonData(c)
	req := getRequest(jsonData)
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
}

func getJsonData(c *comic.Comic) string {
	lastPage := c.Get10LastComicPagesNewestFirst()[10]
	title := lastPage.GetTitle()
	imgurl := lastPage.GetImgUrl()
	prodUrl := lastPage.GetProdUrl()
	description := lastPage.GetDescription()
	tags := "webcomic,graphicnovel,comic,comicart,comics,sciencefiction,scifi,geek,nerd,art,artist,artwork,blackandwhite,concept,conceptart,create,creative,design,digital,draw,drawing,drawings,dystopy,fantasy,humor,illustration,illustrator,image,imagination,ink,inked,inking,kunst,malen,malerei,narrative,parody,pulp,sketch,sketchbook,tusche,zeichnen,zeichnung"
	return fmt.Sprintf(`{"Link":"%s","ImgUrl":"%s","Title":"%s","TagsCsvString":"%s","Description":"%s"}`, prodUrl, imgurl, title, tags, description)
}

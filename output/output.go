package output

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/ingmardrewing/gomic/comic"
	"github.com/ingmardrewing/gomic/config"
	"github.com/ingmardrewing/gomic/page"
)

type Output struct {
	comic  *comic.Comic
	config *config.Config
}

func NewOutput(comic *comic.Comic, config *config.Config) *Output {
	return &Output{comic, config}
}

func (o *Output) WriteToFilesystem() {
	for _, p := range o.comic.GetPages() {
		o.writePageToFileSystem(p)
	}
	o.writeCss()
}

func (o *Output) writeCss() {
	p := o.config.Rootpath + "/css"
	o.prepareFileSystem(p)
	fp := p + "/style.css"
	o.writeStringToFS(fp, css)
}

func (o *Output) writePageToFileSystem(p *page.Page) {
	absPath := o.config.Rootpath + p.FSPath()
	o.prepareFileSystem(absPath)
	o.writeStringToFS(absPath+"/index.html", p.Html())
}

func (o *Output) writeStringToFS(absPath string, html string) {
	log.Println("writing html to filesystem: ", absPath)
	b := []byte(html)
	err := ioutil.WriteFile(absPath, b, 0644)
	if err != nil {
		panic(err)
	}
}

func (o *Output) prepareFileSystem(absPath string) {
	exists, err := o.pathExists(absPath)
	if err != nil {
		panic(err.Error())
	}
	if !exists {
		log.Println("creating path", absPath)
		os.MkdirAll(absPath, 0755)
	}
}

func (o *Output) pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

const css = `
body {
	background-color: red;
}
`

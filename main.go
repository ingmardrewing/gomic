package main

import (
	"github.com/ingmardrewing/gomic/comic"
	"github.com/ingmardrewing/gomic/config"
	"github.com/ingmardrewing/gomic/output"
	"github.com/ingmardrewing/gomic/page"
)

func main() {
	conf := config.NewConfig("/Users/drewing/Sites/gomic.yaml")
	comic := comic.NewComic(conf.Rootpath)

	for _, pc := range conf.Pages {
		p := page.NewPage(pc["title"], pc["path"], pc["imgUrl"], conf.Servedrootpath)
		comic.AddPage(p)
	}

	comic.ConnectPages()

	output := output.NewOutput(&comic, conf)
	output.WriteToFilesystem()
}

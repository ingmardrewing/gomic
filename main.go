package main

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"

	"github.com/ingmardrewing/gomic/comic"
	"github.com/ingmardrewing/gomic/page"
)

func main() {
	config := readConfig()
	comic := comic.NewComic(config.Rootpath)

	for _, pc := range config.Pages {
		p := page.NewPage(pc["title"], pc["path"], pc["imgUrl"])
		comic.AddPage(p)
	}

	comic.ConnectPages()
	comic.PrintPages()
}

type Config struct {
	Url      string              `yaml:"url"`
	Rootpath string              `yaml:"rootpath"`
	Pages    []map[string]string `yaml:"pages"`
}

func readConfig() *Config {
	configpath := "/Users/drewing/Sites/gomic.yaml"
	yamldata, err := ioutil.ReadFile(configpath)
	if err != nil {
		panic(err.Error())
	}

	var config Config
	err = yaml.Unmarshal(yamldata, &config)
	if err != nil {
		panic(err.Error())
	}

	return &config
}

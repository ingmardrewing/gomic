package config

import (
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

func main() {
	fmt.Println("vim-go")
}

type cnf struct {
	Url            string              `yaml:"url"`
	Rootpath       string              `yaml:"rootpath"`
	Servedrootpath string              `yaml:"servedrootpath"`
	PngDir         string              `yaml:"pngdir"`
	Pages          []map[string]string `yaml:"pages"`
}

var conf *cnf

func Read(yamlPath string) {
	conf = newConfig(yamlPath)
}

func Servedrootpath() string {
	srp := conf.Servedrootpath
	log.Println("Servedrootpath", srp)
	return srp
}

func Pages() []map[string]string {
	pgs := conf.Pages
	log.Println(pgs)
	return pgs
}

func Rootpath() string {
	rp := conf.Rootpath
	log.Println(rp)
	return rp
}

func PngDir() string {
	pd := conf.PngDir
	log.Println(pd)
	return pd
}

func newConfig(yamlPath string) *cnf {

	yamldata, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		panic(err.Error())
	}

	var config cnf
	err = yaml.Unmarshal(yamldata, &config)
	if err != nil {
		panic(err.Error())
	}

	return &config
}

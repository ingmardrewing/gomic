package config

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

func main() {
	fmt.Println("vim-go")
}

type Config struct {
	Url      string              `yaml:"url"`
	Rootpath string              `yaml:"rootpath"`
	Pages    []map[string]string `yaml:"pages"`
}

func NewConfig(yamlPath string) *Config {

	yamldata, err := ioutil.ReadFile(yamlPath)
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

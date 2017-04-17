package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

func main() {
	fmt.Println("vim-go")
}

type cnf struct {
	Url                string              `yaml:"url"`
	AwsBucket          string              `yaml:"aws_bucket"`
	AwsDir             string              `yaml:"aws_dir"`
	Rootpath           string              `yaml:"rootpath"`
	Servedrootpath     string              `yaml:"servedrootpath"`
	ServedTestrootpath string              `yaml:"servedtestrootpath"`
	ServedProdrootpath string              `yaml:"servedprodrootpath"`
	PngDir             string              `yaml:"pngdir"`
	Pages              []map[string]string `yaml:"pages"`
}

var conf *cnf
var Stage string

func Read(yamlPath string) {
	conf = newConfig(yamlPath)
}

func IsDev() bool {
	return Stage == "dev"
}

func IsProd() bool {
	return Stage == "prod"
}

func IsTest() bool {
	return Stage == "test"
}

func Servedrootpath() string {
	if IsProd() {
		return conf.ServedProdrootpath
	}
	if IsTest() {
		return conf.ServedTestrootpath
	}
	if IsDev() {
		return conf.Servedrootpath
	}
	return ""
}

func Pages() []map[string]string {
	pgs := conf.Pages
	log.Println(pgs)
	return pgs
}

func Rootpath() string {
	rp := conf.Rootpath
	return rp
}

func PngDir() string {
	pd := conf.PngDir
	return pd
}

func GetTumblData() (string, string, string, string) {
	consumer_key := os.Getenv("GOMIC_TUMBLR_CONSUMER_KEY")
	consumer_secret := os.Getenv("GOMIC_TUMBLR_CONSUMER_SECRET")
	token := os.Getenv("GOMIC_TUMBLR_TOKEN")
	token_secret := os.Getenv("GOMIC_TUMBLR_TOKEN_SECRET")
	return consumer_key, consumer_secret, token, token_secret
}

func GetDsn() string {
	user := os.Getenv("DB_GOMIC_USER")
	pass := os.Getenv("DB_GOMIC_PASS")
	name := os.Getenv("DB_GOMIC_NAME")
	host := os.Getenv("DB_GOMIC_HOST")
	return fmt.Sprintf("%s:%s@%s/%s", user, pass, host, name)
}

func SshUsername() string {
	return os.Getenv("SSH_USER")
}

func SshKeyfilePath() string {
	return os.Getenv("SSH_KEY")
}

func ReadAwsRegion() string {
	return os.Getenv("AWS_REGION")
}

func AwsBucket() string {
	return conf.AwsBucket
}

func AwsDir() string {
	return conf.AwsDir
}

func newConfig(yamlPath string) *cnf {
	stg := flag.String("stage", "", "target stage")
	flag.Parse()
	Stage = *stg
	if Stage == "" {
		fmt.Println(`Usage:

		gomic -stage=<stage>

where <stage> is one of dev, prod, test`)
		os.Exit(0)
	}

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

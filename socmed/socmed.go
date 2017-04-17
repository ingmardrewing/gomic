package socmed

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/MariaTerzieva/gotumblr"
	"github.com/ingmardrewing/gomic/config"
)

var newSocmed = map[string]string{}
var imgurl = ""

func Prepare(p string, t string, i string) {
	newSocmed[p] = t
	imgurl = i
}

func TweetCascade() {
	fmt.Println("tweeting ...")
	for p, t := range newSocmed {
		command := "/Users/drewing/bin/tweetNewComic.pl"
		args := []string{"'" + t + "'", p}
		if err := exec.Command(command, args...).Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	fmt.Println("done")
}

func PublishOnFacebook() {
	// pl gets data from mysql database:
	fmt.Println("Publishing on facebook.")
	command := "open"
	args := []string{"http://localhost/~drewing/cgi-bin/fb.pl"}
	if err := exec.Command(command, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("done")
}

func PostToTumblr() {
	cons_key, cons_secret, token, token_secret := config.GetTumblData()
	fmt.Println(cons_key)
	fmt.Println(cons_secret)
	fmt.Println(token)
	fmt.Println(token_secret)

	client := gotumblr.NewTumblrRestClient(cons_key, cons_secret, token, token_secret, "http://localhost/~drewing/cgi-bin/tumblr.pl", "http://api.tumblr.com")

	blogname := "devabo-de.tumblr.com"
	state := "published"
	photoPostByURL := client.CreatePhoto(blogname, map[string]string{"source": imgurl, "state": state})
	fmt.Println(photoPostByURL)
}

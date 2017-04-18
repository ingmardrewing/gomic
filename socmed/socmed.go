package socmed

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/MariaTerzieva/gotumblr"
	"github.com/ingmardrewing/gomic/config"
)

var imgurl = ""
var prodUrl = ""
var title = ""
var path = ""

func Prepare(p string, t string, i string, pu string) {
	title = t
	path = p
	imgurl = i
	prodUrl = pu
}

func TweetCascade() {
	fmt.Println("tweeting ...")
	command := "/Users/drewing/bin/tweetNewComic.pl"
	args := []string{"'" + title + "'", path}
	if err := exec.Command(command, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
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
	fmt.Println("Post to tumblr")
	cons_key, cons_secret, token, token_secret := config.GetTumblData()
	client := gotumblr.NewTumblrRestClient(cons_key, cons_secret, token, token_secret, "http://localhost/~drewing/cgi-bin/tumblr.pl", "http://api.tumblr.com")

	blogname := "devabo-de.tumblr.com"
	state := "published"
	tags := "comic,webcomic,graphicnovel,drawing,art,narrative,scifi,sci-fi,science-fiction,dystopy,parody,humor,nerd,pulp,geek,blackandwhite"
	photoPostByURL := client.CreatePhoto(
		blogname,
		map[string]string{
			"link":    prodUrl,
			"source":  imgurl,
			"caption": title,
			"tags":    tags,
			"state":   state})
	if photoPostByURL == nil {
		fmt.Println("done")
	} else {
		fmt.Println(photoPostByURL)
	}
}

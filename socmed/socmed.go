package socmed

import (
	"fmt"
	"os"
	"os/exec"
)

var newSocmed = map[string]string{}

func Prepare(p string, t string) {
	newSocmed[p] = t
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

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
	for p, t := range newSocmed {
		command := "/Users/drewing/bin/tweetNewComic.pl"
		args := []string{"'" + t + "'", p}
		if err := exec.Command(command, args...).Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Tweeted %s: %s\n", p, t)
	}
}

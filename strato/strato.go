package strato

import (
	"fmt"
	"os"
	"os/exec"
)

func UploadTest() {
	args := []string{"-r", "/Users/drewing/Sites/gomic", "www.drewing.de@ssh.strato.de:devabo.de/"}
	upload("scp", args)
}

func UploadProd() {
	args := []string{}
	upload("upload_gomic_prod", args)
}

func upload(command string, args []string) {
	if err := exec.Command(command, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("upload to strato complete")
}

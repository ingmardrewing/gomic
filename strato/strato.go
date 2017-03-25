package strato

import (
	"fmt"
	"os"
	"os/exec"
)

func UploadDir(path string) {
	cmd := "scp"
	args := []string{"-r", path, "www.drewing.de@ssh.strato.de:devabo.de/"}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("upload to strato complete")
}

package pkg

import (
	"bytes"
	"os/exec"
	"strings"
)

// https://unix.stackexchange.com/questions/247576/how-to-get-home-given-user
func getHomeDirectory() (home string, err error) {
	whoami, err := run(exec.Command("whoami"))
	if err != nil {
		return
	}
	getent, err := run(exec.Command("getent", "passwd", whoami))
	if err != nil {
		return
	}
	cut := exec.Command("cut", "-d:", "-f6")
	cut.Stdin = bytes.NewReader([]byte(getent))
	return run(cut)
}

func run(cmd *exec.Cmd) (string, error) {
	output, err := cmd.Output()
	return strings.TrimSpace(string(output)), err
}

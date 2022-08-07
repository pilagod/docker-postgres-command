package pkg

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func createPgPass(conn Connection, path string) error {
	if err := os.WriteFile(
		path,
		// hostname:port:database:username:password
		[]byte(fmt.Sprintf("%s:%s:%s:%s:%s", conn.Host, conn.Port, conn.DB, conn.Username, conn.Password)),
		0o600,
	); err != nil {
		return fmt.Errorf("Create %s fails: %s", path, err.Error())
	}
	return nil
}

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

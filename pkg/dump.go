package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type DumpOption struct {
	Connection
	Flags string
}

func Dump(opt DumpOption) (path string, err error) {
	home, err := getHomeDirectory()
	if err != nil {
		err = fmt.Errorf("Cannot get home directory: %s", err.Error())
		return
	}

	// Dump db to home directory
	path = fmt.Sprintf(
		"%s/pgdb_%s.dump",
		home,
		time.Now().Format("20060102150405"),
	)

	// Create .pgpass file at home directory
	pgPassPath := home + "/.pgpass"
	if err = os.WriteFile(
		pgPassPath,
		// hostname:port:database:username:password
		[]byte(fmt.Sprintf("%s:%s:%s:%s:%s", opt.Host, opt.Port, opt.DB, opt.Username, opt.Password)),
		0o600,
	); err != nil {
		err = fmt.Errorf("Create %s fails: %s", pgPassPath, err.Error())
		return
	}
	defer func() {
		if err = os.Remove(pgPassPath); err != nil {
			err = fmt.Errorf("Remove %s fails: %s", pgPassPath, err.Error())
			return
		}
	}()

	// First element after splitting is empty string, just omit it.
	pgDumpArgs := strings.Split(
		fmt.Sprintf("%s -f%s -d%s -h%s -p%s -U%s -w", opt.Flags, path, opt.DB, opt.Host, opt.Port, opt.Username),
		"-",
	)[1:]
	for i, arg := range pgDumpArgs {
		pgDumpArgs[i] = "-" + strings.TrimSpace(arg)
	}
	if err = exec.Command("pg_dump", pgDumpArgs...).Run(); err != nil {
		err = fmt.Errorf("Execute pg_dump fails: %s", err.Error())
	}
	return
}

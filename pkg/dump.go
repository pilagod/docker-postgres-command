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
	// Dump db to root directory
	path = fmt.Sprintf(
		"/pgdb_%s.dump",
		time.Now().Format("20060102150405"),
	)

	// Create .pgpass file at root directory
	pgPassPath := "/.pgpass"
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
		}
	}()

	// First element after splitting is empty string, just omit it.
	pgDumpArgs := strings.Split(
		fmt.Sprintf("%s -f%s -h%s -p%s -d%s -U%s -w", opt.Flags, path, opt.Host, opt.Port, opt.DB, opt.Username),
		"-",
	)[1:]
	for i, arg := range pgDumpArgs {
		pgDumpArgs[i] = "-" + strings.TrimSpace(arg)
	}
	pgDump := exec.Command("pg_dump", pgDumpArgs...)
	pgDump.Env = append(pgDump.Env, "PGPASSFILE=/.pgpass")
	if err = pgDump.Run(); err != nil {
		err = fmt.Errorf("Execute pg_dump fails: %s", err.Error())
	}
	return
}

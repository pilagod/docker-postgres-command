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
	if err = createPgPass(opt.Connection, pgPassPath); err != nil {
		err = fmt.Errorf("Cannot create %s: %v", pgPassPath, err)
		return
	}
	defer os.Remove(pgPassPath)

	// First element after splitting is empty string, just omit it.
	pgDumpArgs := strings.Split(
		fmt.Sprintf("%s -f%s -h%s -p%s -d%s -U%s -w", opt.Flags, path, opt.Host, opt.Port, opt.DB, opt.Username),
		" ",
	)
	pgDump := exec.Command("pg_dump", pgDumpArgs...)
	pgDump.Env = append(pgDump.Env, "PGPASSFILE=/.pgpass")
	if err = pgDump.Run(); err != nil {
		err = fmt.Errorf("Execute pg_dump fails: %v", err)
	}
	return
}

package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type RestoreOption struct {
	Connection
	Path  string
	Flags string
}

func Restore(opt RestoreOption) error {
	pgPassPath := "/.pgpass"
	if err := createPgPass(opt.Connection, pgPassPath); err != nil {
		return fmt.Errorf("Cannot create %s: %v", pgPassPath, err)
	}
	defer os.Remove(pgPassPath)

	pgRestoreArgs := strings.Split(
		fmt.Sprintf("%s -h%s -p%s -d%s -U%s -w %s", opt.Flags, opt.Host, opt.Port, opt.DB, opt.Username, opt.Path),
		" ",
	)
	pgRestore := exec.Command("pg_restore", pgRestoreArgs...)
	pgRestore.Env = append(pgRestore.Env, fmt.Sprintf("PGPASSFILE=%s", pgPassPath))
	if err := pgRestore.Run(); err != nil {
		return fmt.Errorf("Execute pg_restore fails: %v", err)
	}
	return nil
}

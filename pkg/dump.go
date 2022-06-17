package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Dump() (err error) {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	db := os.Getenv("DB")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	file := os.Getenv("DUMP_FILE")
	flags := os.Getenv("DUMP_FLAGS")

	// Create .pgpass file at home directory
	home, err := getHomeDirectory()
	if err != nil {
		return fmt.Errorf("Cannot get home directory: %s", err.Error())
	}
	pgPassPath := home + "/.pgpass"
	if err = os.WriteFile(
		pgPassPath,
		// hostname:port:database:username:password
		[]byte(fmt.Sprintf("%s:%s:%s:%s:%s", host, port, db, username, password)),
		0o600,
	); err != nil {
		return fmt.Errorf("Create %s fails: %s", pgPassPath, err.Error())
	}

	// First element after splitting is empty string, just omit it.
	pgDumpArgs := strings.Split(
		fmt.Sprintf("%s -f%s -d%s -h%s -p%s -U%s -w", flags, file, db, host, port, username),
		"-",
	)[1:]
	for i, arg := range pgDumpArgs {
		pgDumpArgs[i] = "-" + strings.TrimSpace(arg)
	}
	if err = exec.Command("pg_dump", pgDumpArgs...).Run(); err != nil {
		return fmt.Errorf("Execute pg_dump fails: %s", err.Error())
	}

	// Remove .pgpass file
	if err = os.Remove(pgPassPath); err != nil {
		return fmt.Errorf("Remove %s fails: %s", pgPassPath, err.Error())
	}
	return
}

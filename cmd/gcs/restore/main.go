package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"

	"github.com/pilagod/docker-postgres-command/pkg"
)

func main() {
	bucket := os.Getenv("BUCKET")
	object := os.Getenv("OBJECT")

	pgDumpPath := "/" + object
	pgDumpFile, err := os.Create(pgDumpPath)
	if err != nil {
		panic(fmt.Errorf("os.Create: %v", err))
	}
	defer os.Remove(pgDumpPath)

	download(bucket, object, pgDumpFile)

	fmt.Printf("Blob %v downloaded\n", object)

	if err := pkg.Restore(pkg.RestoreOption{
		Connection: pkg.Connection{
			Host:     os.Getenv("HOST"),
			Port:     os.Getenv("PORT"),
			DB:       os.Getenv("DB"),
			Username: os.Getenv("USERNAME"),
			Password: os.Getenv("PASSWORD"),
		},
		Path:  pgDumpPath,
		Flags: os.Getenv("RESTORE_FLAGS"),
	}); err != nil {
		panic(fmt.Errorf("Restore %s fails: %v", pgDumpPath, err))
	}

	fmt.Printf("Restore %s\n", pgDumpPath)
}

func download(bucket, object string, f io.WriteCloser) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(fmt.Errorf("storage.NewClient: %v", err))
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		panic(fmt.Errorf("Object(%q).NewReader: %v", object, err))
	}
	defer rc.Close()

	if _, err := io.Copy(f, rc); err != nil {
		panic(fmt.Errorf("io.Copy: %v", err))
	}
	if err := f.Close(); err != nil {
		panic(fmt.Errorf("f.Close: %v", err))
	}
}

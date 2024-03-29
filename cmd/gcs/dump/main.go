package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"

	"github.com/pilagod/docker-postgres-command/pkg"
)

func main() {
	pgDumpPath, err := pkg.Dump(
		pkg.DumpOption{
			Connection: pkg.Connection{
				Host:     os.Getenv("HOST"),
				Port:     os.Getenv("PORT"),
				DB:       os.Getenv("DB"),
				Username: os.Getenv("USERNAME"),
				Password: os.Getenv("PASSWORD"),
			},
			Flags: os.Getenv("DUMP_FLAGS"),
		},
	)
	if err != nil {
		panic(err)
	}
	defer os.Remove(pgDumpPath)

	bucket := os.Getenv("BUCKET")
	object := pgDumpPath[strings.LastIndex(pgDumpPath, "/")+1:]

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(fmt.Errorf("storage.NewClient: %v", err))
	}
	defer client.Close()

	pgDumpFile, err := os.Open(pgDumpPath)
	if err != nil {
		panic(fmt.Errorf("os.Open: %v", err))
	}
	defer pgDumpFile.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	o := client.Bucket(bucket).Object(object)
	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	// For an object that does not yet exist, set the DoesNotExist precondition.
	o = o.If(storage.Conditions{DoesNotExist: true})

	// If the live object already exists in your bucket, set instead a
	// generation-match precondition using the live object's generation number.
	// attrs, err := o.Attrs(ctx)
	// if err != nil {
	// 	return fmt.Errorf("object.Attrs: %v", err)
	// }
	// o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	// Upload an object with storage.Writer.
	wc := o.NewWriter(ctx)
	if _, err := io.Copy(wc, pgDumpFile); err != nil {
		panic(fmt.Errorf("io.Copy: %v", err))
	}
	if err := wc.Close(); err != nil {
		panic(fmt.Errorf("Writer.Close: %v", err))
	}
	fmt.Printf("Blob %v uploaded.\n", object)
}

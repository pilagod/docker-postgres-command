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
	path, err := pkg.Dump(
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
	bucket := os.Getenv("BUCKET")
	object := path[strings.LastIndex(path, "/")+1:]

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(fmt.Errorf("storage.NewClient: %v", err))
	}
	defer client.Close()

	// Open local file.
	f, err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("os.Open: %v", err))
	}
	defer f.Close()

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
	if _, err := io.Copy(wc, f); err != nil {
		panic(fmt.Errorf("io.Copy: %v", err))
	}
	if err := wc.Close(); err != nil {
		panic(fmt.Errorf("Writer.Close: %v", err))
	}
	fmt.Printf("Blob %v uploaded.\n", object)

	if err := os.Remove(path); err != nil {
		panic(fmt.Errorf("Cannot remove dump file %s: %v", path, err))
	}
}

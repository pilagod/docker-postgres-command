package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

func main() {
	bucket := os.Getenv("BUCKET")
	object := os.Getenv("OBJECT")

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(fmt.Errorf("storage.NewClient: %v", err))
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	f, err := os.Create("/" + object)
	if err != nil {
		panic(fmt.Errorf("os.Create: %v", err))
	}

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

	fmt.Printf("Blob %v downloaded\n", object)
}

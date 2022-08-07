build-gcs-dump:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./bin/dump ./cmd/gcs/dump/main.go

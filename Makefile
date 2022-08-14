build-docker:
	make build-gcs-dump
	make build-gcs-restore
	docker build --build-arg POSTGRES=$(POSTGRES) -t $(REGISTRY)/postgres-command:$(POSTGRES) .

build-gcs-dump:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./bin/gcs-dump ./cmd/gcs/dump/main.go

build-gcs-restore:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./bin/gcs-restore ./cmd/gcs/restore/main.go

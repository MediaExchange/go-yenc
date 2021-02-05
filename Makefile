.PHONY: build test

build:
	go build -o yenc ./cmd/main.go

test:
	go test ./...

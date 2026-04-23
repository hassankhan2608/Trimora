.PHONY: run build test vet fmt lint tidy

run:
	go run ./cmd/server

build:
	go build -o bin/trimora ./cmd/server

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

lint:
	golangci-lint run

tidy:
	go mod tidy

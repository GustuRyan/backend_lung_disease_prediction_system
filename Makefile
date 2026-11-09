.PHONY: run build migrate test fmt vet

run:
	go run ./cmd/api

build:
	go build ./cmd/api

migrate:
	go run ./cmd/migrate

test:
	go test -race ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

vet:
	go vet ./...

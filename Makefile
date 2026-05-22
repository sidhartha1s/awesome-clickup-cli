.PHONY: build test lint install clean

build:
	go build -o bin/clickup-reference-pp-cli ./cmd/clickup-reference-pp-cli

test:
	go test ./...

lint:
	golangci-lint run

install:
	go install ./cmd/clickup-reference-pp-cli

clean:
	rm -rf bin/

build-mcp:
	go build -o bin/clickup-reference-pp-mcp ./cmd/clickup-reference-pp-mcp

install-mcp:
	go install ./cmd/clickup-reference-pp-mcp

build-all: build build-mcp

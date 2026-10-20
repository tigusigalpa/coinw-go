.PHONY: test test-verbose vet fmt build clean help

all: test

test:
	@go test ./...

test-verbose:
	@go test -v ./...

vet:
	@go vet ./...

fmt:
	@gofmt -w .

build:
	@go build ./...

clean:
	@rm -f coverage.out

help:
	@echo "Targets: test test-verbose vet fmt build clean"

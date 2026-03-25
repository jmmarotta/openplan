.PHONY: all vet test build

all: vet test build

vet:
	go vet ./...

test:
	go test ./...

build:
	go build ./...

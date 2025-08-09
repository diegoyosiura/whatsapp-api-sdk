.PHONY: all build lint test clean

all: build

build:
	go build -o bin/whatsapp-sdk ./...

lint:
	golangci-lint run ./...

test:
	go test ./... -cover

clean:
	rm -rf bin/ build/ dist/
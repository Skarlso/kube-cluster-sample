BINARY=receiver

.DEFAULT_GOAL := build

.PHONY: clean build test


build:
	go build -i -o ./build/${BINARY} && cp .env build/

test:
	go test ./...

get-deps:
	dep ensure

clean:
	go clean -i

run:
	./build/${BINARY}

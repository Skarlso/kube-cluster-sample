BINARY=frontend

.DEFAULT_GOAL := build

.PHONY: clean build test


build:
	go build -i -o ./build/${BINARY}

linux:
	env GOOS=linux GOARCH=arm go build -o build/${BINARY}

test:
	go test ./...

get-deps:
	dep ensure

clean:
	go clean -i

run:
	./build/${BINARY}

BINARY=frontend

# Set the build dir, where built cross-compiled binaries will be output
BUILDDIR := build
# List the GOOS and GOARCH to build
GO_LDFLAGS_STATIC="-s -w $(CTIMEVAR) -extldflags -static"

.DEFAULT_GOAL := binaries

binaries:
	CGO_ENABLED=0 gox \
		-osarch="linux/amd64 linux/arm darwin/amd64" \
		-ldflags=${GO_LDFLAGS_STATIC} \
		-output="$(BUILDDIR)/{{.OS}}/{{.Arch}}/$(BINARY)" \
		-tags="netgo" \
		./...

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

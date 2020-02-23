NAME=receiver

# Set the build dir, where built cross-compiled binaries will be output
BUILDDIR := build

# List the GOOS and GOARCH to build
GO_LDFLAGS_STATIC="-s -w $(CTIMEVAR) -extldflags -static"

.DEFAULT_GOAL := binaries

binaries:
	CGO_ENABLED=0 gox \
		-osarch="linux/amd64 linux/arm darwin/amd64" \
		-ldflags=${GO_LDFLAGS_STATIC} \
		-output="$(BUILDDIR)/{{.OS}}/{{.Arch}}/$(NAME)" \
		-tags="netgo" \
		./...

.PHONY: build
build:
	go build -ldflags="-s -w" -i -o ${BUILDDIR}/${NAME} cmd/root.go

.PHONY: bootstrap
bootstrap:
	go get github.com/mitchellh/gox

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	go clean -i

lint:
	golint -set_exit_status ./...

.PHONY: run
run:
	go run cmd/root.go

docker_image:
	docker build -t $(image):$(version) .
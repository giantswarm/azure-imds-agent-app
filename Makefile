PROJECT=azure-imds-agent-app
ORGANISATION=giantswarm
VERSION=0.1.0
BIN=$(PROJECT)
GOVERSION := 1.15.2
BUILDDATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
COMMIT := $(shell git rev-parse HEAD)
SOURCE=$(shell find . -name '*.go')

ifndef GOOS
	GOOS := $(shell go env GOOS)
endif
ifndef GOARCH
	GOARCH := $(shell go env GOARCH)
endif

# binary to test with
TESTBIN := build/bin/${BIN}-${GOOS}-${GOARCH}

.PHONY: all azure-imds-agent-app

all: build docker-build

fetch:
	docker pull golang:$(GOVERSION)-alpine

# build binary for current platform
build: fetch build/bin/$(BIN)-$(GOOS)-$(GOARCH)

azure-imds-agent-app:
	@go build \
		-o azure-imds-agent-app \
		-ldflags " \
			-X main.gitCommit=$(COMMIT) \
		" \
		.

# platform-specific build for linux-amd64
# - here we have CGO_ENABLED=1
build/bin/$(BIN)-linux-amd64: $(SOURCE)
	@mkdir -p build/bin
	docker run --rm -v $(shell pwd):/go/src/github.com/$(ORGANISATION)/$(PROJECT) \
		-e GOPATH=/go -e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=1 \
		-w /go/src/github.com/$(ORGANISATION)/$(PROJECT) \
		quay.io/giantswarm/golang:$(GOVERSION) go build -a -installsuffix cgo -o build/bin/$(BIN)-linux-amd64 \
		-ldflags "-X 'main.gitCommit=$(COMMIT)'"

docker-build: build
	docker build --no-cache -t quay.io/$(ORGANISATION)/$(PROJECT):latest .
	docker tag quay.io/$(ORGANISATION)/$(PROJECT):latest quay.io/$(ORGANISATION)/$(PROJECT):$(VERSION)

# run unittests
gotest:
	go test -cover ./...
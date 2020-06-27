# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=kubectl-recreate
BINARY_UNIX=$(BINARY_NAME)
BINARY_WINDOWS=$(BINARY_NAME).exe

VERSION?=?
RELEASE_VERSION=$(shell git describe --tags $(git rev-list --tags --max-count=1))
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"
RELEASE_LDFLAGS = -ldflags "-X main.VERSION=${RELEASE_VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

all: build
lint:
	wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.23.7
	./bin/golangci-lint run ./...		
test:
	$(GOTEST) -v ./...
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ${LDFLAGS}
install:
	$(GOBUILD) -o $(BINARY_NAME) -v ${LDFLAGS}
	mv $(BINARY_NAME)  /usr/local/bin
e2e-test:
	$(GOTEST) -v  ./test/e2e



# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
build-windows:
	GOOS=windows GOARCH=386 $(GOBUILD) -o $(BINARY_WINDOWS) -v
build-release:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o release/$(BINARY_UNIX) -v ${RELEASE_LDFLAGS}
	GOOS=windows GOARCH=386 $(GOBUILD) -o release/$(BINARY_WINDOWS) -v ${RELEASE_LDFLAGS}

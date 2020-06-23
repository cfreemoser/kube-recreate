# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=kubectl-refresh
BINARY_UNIX=$(BINARY_NAME)
BINARY_WINDOWS=$(BINARY_NAME).exe

all: build
test:
	$(GOTEST) -v ./...
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
install:
	$(GOBUILD) -o $(BINARY_NAME) -v
	mv $(BINARY_NAME)  /usr/local/bin


# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
build-windows:
	GOOS=windows GOARCH=386 $(GOBUILD) -o $(BINARY_WINDOWS) -v
build-release:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o release/$(BINARY_UNIX) -v
	GOOS=windows GOARCH=386 $(GOBUILD) -o release/$(BINARY_WINDOWS) -v

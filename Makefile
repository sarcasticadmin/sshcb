# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY_NAME=sshcb

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

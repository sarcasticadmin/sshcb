# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=sshcb
TEST_DIRS = builder

PKG := github.com/sarcasticadmin/sshcb
VERSION := $(shell git describe --dirty)
LDFLAGS := -X $(PKG)/cmd.Version=$(VERSION)

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v -ldflags="$(LDFLAGS)"
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
test:
	$(GOTEST) -v -cover ./builder

test-release:
	goreleaser release --rm-dist --skip-publish --skip-validate

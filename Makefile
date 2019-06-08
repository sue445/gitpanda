# Requirements: git, go, vgo
NAME     := gitpanda
VERSION  := $(shell cat VERSION)
REVISION := $(shell git rev-parse --short HEAD)

SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): $(SRCS)
	GO111MODULE=on go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: clean
clean:
	rm -rf bin/*

.PHONY: tag
tag:
	git tag -a $(VERSION) -m "Release v$(VERSION)"
	git push --tags

.PHONY: release
release: tag
	git push origin master

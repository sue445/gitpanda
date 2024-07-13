# Requirements: git, go, vgo
NAME     := gitpanda
VERSION  := $(shell cat VERSION)
REVISION := $(shell git rev-parse --short HEAD)

SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := "-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): $(SRCS)
	GO111MODULE=on go build -ldflags=$(LDFLAGS) -o bin/$(NAME)

.PHONY: linux_amd64
linux_amd64: $(SRCS)
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -ldflags=$(LDFLAGS) -o bin/$(NAME)_linux_amd64

.PHONY: gox
gox:
	gox -osarch="$${GOX_OSARCH}" -ldflags=$(LDFLAGS) -output="bin/gitpanda_{{.OS}}_{{.Arch}}"

.PHONY: zip
zip:
	cd bin ; \
	for file in *; do \
		zip $${file}.zip $${file} ; \
	done

.PHONY: gox_with_zip
gox_with_zip: clean gox zip

.PHONY: clean
clean:
	rm -rf bin/*

.PHONY: tag
tag:
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push --tags

.PHONY: release
release: tag
	git push origin main

.PHONY: test
test:
	go test -count=1 $${TEST_ARGS} ./...

.PHONY: testrace
testrace:
	go test -count=1 $${TEST_ARGS} -race ./...

.PHONY: fmt
fmt:
	go fmt ./...

SHELL := /bin/bash

PREFIX = sshare

OS=$(or ${GOOS},${GOOS},linux)
ARCH=$(or ${GOARCH},${GOARCH},amd64)

CURRENTDIR = $(shell pwd)
SOURCEDIR = $(CURRENTDIR)

PATH := $(CURRENTDIR)/bin:$(PATH)

VERSION?=$(shell git describe --tags --always)

LD_FLAGS = -ldflags "-X 'sshare/pkg/version.VERSION=$(VERSION)' -s -w"

.PHONY: clean build docker-build docker-push

build: dist/sshare

dist/sshare:
	mkdir -p $(@D)
	GOOS=${OS} GOARCH=${GOARCH} go build $(LD_FLAGS) -v -o $(@D)

clean:
	rm -rf dist

docker-build: build
	docker build -t tczekajlo/$(PREFIX):$(VERSION) .

docker-push:
	docker push tczekajlo/$(PREFIX):$(VERSION)

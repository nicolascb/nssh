SHELL := $(shell which bash)
BINARY := "nssh"
OSARCH := "linux/amd64 linux/386 darwin/amd64 darwin/386"
ENV = /usr/bin/env

.SHELLFLAGS = -c
.SILENT: ;
.ONESHELL: ;
.NOTPARALLEL: ;
.EXPORT_ALL_VARIABLES: ;

.PHONY: all
.DEFAULT: help

help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

dep: ## Get build dependencies
	go get -v -u github.com/golang/dep/cmd/dep && \
	go get github.com/mitchellh/gox

build: ## Build the app
	dep ensure && go build -o bin/$(BINARY)

cross-build: ## Build the app for multiple os/arch
	gox -osarch=$(OSARCH) -output "bin/$(BINARY)_{{.OS}}_{{.Arch}}"

test: ## Launch tests
	go test -v *.go
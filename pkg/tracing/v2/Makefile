UNAME_S = $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	SHELL=/usr/local/bin/bash
else
	SHELL=/bin/bash
endif
.SHELLFLAGS = -o pipefail -O globstar -c

PROJECT_NAME = $(shell basename `git rev-parse --show-toplevel`)
ROOT_PATH = $(abspath . )
PKG = $(ROOT_PATH)/$(PROJECT_NAME)
PKG_LIST = $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: gen
gen:
	go generate ./...

.PHONY: fmt
fmt:
	fieldalignment -fix ./... || true
	gofmt -w $$(gofmt -l ./ | xargs) || true

.PHONY: lint
lint:
	golangci-lint run $(PKG_LIST)

.PHONY: test
test:
	mkdir -p /tmp/$(PROJECT_NAME) && \
	go test $(PKG_LIST) -covermode=count -coverprofile=overalls.coverprofile -p 2 && \
        go tool cover -func=./overalls.coverprofile && \
        go tool cover -html=./overalls.coverprofile -o /tmp/$(PROJECT_NAME)/coverage.html && \
        go tool cover -func=./overalls.coverprofile -o /tmp/$(PROJECT_NAME)/coverage.txt && \
        rm -f ./overalls.coverprofile

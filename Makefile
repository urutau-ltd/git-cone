SHELL := bash

BINARY := cone
CMD_DIR := ./cmd/cone
IMAGE ?= ghcr.io/urutau-ltd/git-cone:latest
CONTAINER ?= docker
GO ?= go
GOFLAGS ?=
GOCACHE ?= $(CURDIR)/.gocache
CGO_ENABLED ?= 0

.PHONY: help shell build build-local install test test-all fmt lint vuln image run clean

help:
	@printf '%s\n' \
		'make build      Build ./cmd/cone with CGO disabled' \
		'make install    Install cone into $$GOBIN or $$GOPATH/bin' \
		'make test       Run the core package test suite used by the build contract' \
		'make test-all   Run the expanded package suite used for this fork' \
		'make fmt        Run gofmt on the repository' \
		'make lint       Run staticcheck' \
		'make vuln       Run govulncheck' \
		'make image      Build the container image' \
		'make run        Run cone serve locally' \
		'make shell      Enter the Guix dev shell from manifest.scm' \
		'make clean      Remove build artifacts'

shell:
	guix shell -m manifest.scm

build:
	mkdir -p ./dist
	CGO_ENABLED=$(CGO_ENABLED) GOCACHE=$(GOCACHE) $(GO) build $(GOFLAGS) -trimpath -ldflags="-s -w" -o ./dist/$(BINARY) $(CMD_DIR)

build-local:
	mkdir -p ./dist
	CGO_ENABLED=$(CGO_ENABLED) GOCACHE=$(GOCACHE) $(GO) build $(GOFLAGS) -o ./dist/$(BINARY) $(CMD_DIR)

install:
	CGO_ENABLED=$(CGO_ENABLED) GOCACHE=$(GOCACHE) $(GO) install $(GOFLAGS) $(CMD_DIR)

test:
	CGO_ENABLED=$(CGO_ENABLED) GOCACHE=$(GOCACHE) $(GO) test $(GOFLAGS) ./pkg/access/... ./pkg/utils/... ./pkg/webhook/... ./pkg/ssh/...

test-all:
	CGO_ENABLED=$(CGO_ENABLED) GOCACHE=$(GOCACHE) $(GO) test $(GOFLAGS) ./pkg/access/... ./pkg/utils/... ./pkg/webhook/... ./pkg/ssh/... ./pkg/notify/... ./pkg/config/...

fmt:
	find . -name '*.go' -print0 | xargs -0 gofmt -w

lint:
	GOCACHE=$(GOCACHE) staticcheck ./...

vuln:
	GOCACHE=$(GOCACHE) govulncheck ./...

image:
	$(CONTAINER) build -t $(IMAGE) .

run:
	CGO_ENABLED=$(CGO_ENABLED) GOCACHE=$(GOCACHE) $(GO) run $(GOFLAGS) $(CMD_DIR) serve

clean:
	rm -rf ./dist $(GOCACHE)

SHELL := /bin/bash

ifndef VERSION
VERSION := `jq .version package.json `
endif
ifndef GITCOMMIT
GITCOMMIT := `git rev-parse HEAD`
endif
ifndef BUILDDATE
BUILDDATE := `date +%Y-%m-%d`
endif

LDFLAGSSTRING :=-X github.com/ethereum-optimism/op-exporter/version.Version=$(VERSION)
LDFLAGSSTRING +=-X github.com/ethereum-optimism/op-exporter/version.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X github.com/ethereum-optimism/op-exporter/version.BuildDate=$(BUILDDATE)

LDFLAGS :=-ldflags "$(LDFLAGSSTRING)"

.PHONY: all build

all: build

# Build binary
build:
	CGO_ENABLED=0 go build $(LDFLAGS)

lint:
	golangci-lint run -E goimports,sqlclosecheck,bodyclose,asciicheck,misspell,errorlint --timeout 5m -e "errors.As" -e "errors.Is" ./...

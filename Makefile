# Use bash syntax
SHELL=/bin/bash
# Go parameters
GOCMD=go
GOBINPATH=$(shell $(GOCMD) env GOPATH)/bin
GOMOD=$(GOCMD) mod
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=gotestsum
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install
GOTOOL=$(GOCMD) tool
GOFMT=$(GOCMD) fmt

.PHONY: FORCE

.PHONY: all
all: fmt lint test build go.mod

.PHONY: build
build:
	go generate ./...
	go build -v .

.PHONY: test
test:
	$(GOTEST) --format testname --junitfile unit-tests.xml -- -mod=readonly -race -coverprofile=./c.out -covermode=atomic -coverpkg=.,./... ./...

.PHONY: coverage
coverage: test
	$(GOTOOL) cover -func=cover.out

.PHONY: mocks
mocks:
	rm -rf mocks
	mockery --all

.PHONY: fmt
fmt:
	$(GOFMT) ./...

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f bin/*

.PHONY: lint
lint:
	@which golangci-lint > /dev/null 2>&1 || (curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- -b $(GOBINPATH) v1.46.2)
	golangci-lint run -v --timeout=10m -D errcheck

.PHONY: install-tools
install-tools:
	$(GOINSTALL) gotest.tools/gotestsum@v1.6.3
	$(GOINSTALL) github.com/vektra/mockery/v2@latest

go.mod: FORCE
	$(GOMOD) tidy
	$(GOMOD) verify
go.sum: go.mod

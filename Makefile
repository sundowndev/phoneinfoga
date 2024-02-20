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
GIT_TAG=$(shell git describe --abbrev=0 --tags)
GIT_COMMIT=$(shell git rev-parse --short HEAD)

.PHONY: FORCE

.PHONY: all
all: fmt lint test build go.mod

.PHONY: build
build:
	go generate ./...
	go build -v -ldflags="-s -w -X 'github.com/sundowndev/phoneinfoga/v2/build.Version=${GIT_TAG}' -X 'github.com/sundowndev/phoneinfoga/v2/build.Commit=${GIT_COMMIT}'" -o ./bin/phoneinfoga .

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
	golangci-lint run -v --timeout=2m

.PHONY: install-tools
install-tools:
	$(GOINSTALL) gotest.tools/gotestsum@v1.6.3
	$(GOINSTALL) github.com/vektra/mockery/v2@v2.38.0
	$(GOINSTALL) github.com/swaggo/swag/cmd/swag@v1.16.3
	@which golangci-lint > /dev/null 2>&1 || (curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- -b $(GOBINPATH) v1.46.2)

go.mod: FORCE
	$(GOMOD) tidy
	$(GOMOD) verify
go.sum: go.mod

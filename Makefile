BINPATH ?= build

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)

LDFLAGS = -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)"

.PHONY: all
all: audit test build

.PHONY: audit
audit:
	go list -m all | nancy sleuth

.PHONY: lint
lint:
	exit

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: build
build: fmt
	go build -tags 'production' $(LDFLAGS) -o $(BINPATH)/dp-search-data-importer

.PHONY: debug
debug:
	go build -tags 'debug' $(LDFLAGS) -o $(BINPATH)/dp-search-data-importer
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-search-data-importer

.PHONY: test
test: fmt
	go test -count=1 -race -cover ./...

.PHONY: produce
produce:
	HUMAN_LOG=1 go run cmd/producer/main.go

.PHONY: convey
convey:
	goconvey ./...

.PHONY: test-component
test-component:
	go test -cover -coverpkg=github.com/ONSdigital/dp-search-data-importer/... -component

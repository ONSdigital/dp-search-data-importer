#!/bin/bash -eux

pushd dp-search-data-importer
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5
  make lint
  npm install -g @asyncapi/cli
  make validate-specification
popd

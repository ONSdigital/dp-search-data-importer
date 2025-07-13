#!/bin/bash -eux

pushd dp-search-data-importer
  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.2.1
  make lint
  npm install -g @asyncapi/cli
  make validate-specification
popd

#!/bin/bash -eux

pushd dp-search-data-importer
  make lint
  make validate-specification
popd

#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-search-data-importer
  make test-component
popd

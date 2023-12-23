#!/bin/bash

pushd cocotola-api
golangci-lint run --config ../.github/.golangci.yml && \
golangci-lint run --disable-all --config ../.github/.golangci.yml \
-E bodyclose \
-E errorlint \
-E exhaustive \
-E forbidigo \
-E forcetypeassert \
-E gocognit \
-E gocyclo \
-E gofmt \
-E goimports \
-E gomnd \
-E gosec \
-E noctx \
-E unconvert \
-E whitespace && \
pkgforbid -config=../.github/pkgforbid.yml ./... && \
popd

pushd cocotola-auth
golangci-lint run --config ../.github/.golangci.yml && \
golangci-lint run --disable-all --config ../.github/.golangci.yml \
-E bodyclose \
-E errorlint \
-E exhaustive \
-E forbidigo \
-E forcetypeassert \
-E gocognit \
-E gocyclo \
-E gofmt \
-E goimports \
-E gomnd \
-E gosec \
-E noctx \
-E unconvert \
-E whitespace && \
pkgforbid -config=../.github/pkgforbid.yml ./... && \
popd

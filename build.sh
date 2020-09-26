#!/bin/sh

GIT_COMMIT=$(git rev-parse HEAD)
GIT_TAG=$(git describe --tags --exact-match $GIT_COMMIT)

echo "Building linux x64 binary..."
mkdir -p bin/linux/x64
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.gitCommit=$GIT_COMMIT -X main.gitTag=$GIT_TAG" -o bin/linux/x64/galadh cmd/main.go

echo "Building windows x64 binary..."
mkdir -p bin/windows/x64
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.gitCommit=$GIT_COMMIT -X main.gitTag=$GIT_TAG" -o bin/windows/x64/galadh.exe cmd/main.go

echo "Building macos x64 binary..."
mkdir -p bin/macos/x64
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.gitCommit=$GIT_COMMIT -X main.gitTag=$GIT_TAG" -o bin/macos/x64/galadh cmd/main.go
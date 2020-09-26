#!/bin/sh

GIT_COMMIT=$(git rev-parse HEAD)
GIT_TAG=$(git describe --tags --exact-match $GIT_COMMIT)

mkdir -p bin
echo "Building linux x64 binary..."
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.gitCommit=$GIT_COMMIT -X main.gitTag=$GIT_TAG" -o bin/galadh_linux_x64 cmd/main.go

echo "Building windows x64 binary..."
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.gitCommit=$GIT_COMMIT -X main.gitTag=$GIT_TAG" -o bin/galadh_win_x64.exe cmd/main.go

echo "Building macos x64 binary..."
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.gitCommit=$GIT_COMMIT -X main.gitTag=$GIT_TAG" -o bin/galadh_macos_x64 cmd/main.go
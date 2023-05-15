#!/bin/bash

docker run --platform=linux/arm64 --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp -e CGO_ENABLE=0 -e GOOS=linux -e GOARCH=amd64 golang:latest go build -o xray -trimpath -ldflags "-s -w -buildid=" ./main

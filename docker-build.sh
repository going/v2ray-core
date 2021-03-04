#!/bin/bash

docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp -e CGO_ENABLE=0 golang:latest go build -o xray -trimpath -ldflags "-s -w -buildid=" ./main

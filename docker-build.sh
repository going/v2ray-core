#!/bin/bash

docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp -e CGO_ENABLE=0 golang:latest go build -tags "full" -ldflags="-s -w" -v

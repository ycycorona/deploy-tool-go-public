#!/bin/bash -e
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -x -o ./dist/zip-tool-linux ./src/main.go
#!/bin/bash -e
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -x -o ./dist/zip-tool-win.exe ./src/main.go
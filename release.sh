#!/bin/sh

mkdir -p dist
GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.Version=$(git describe)'" -o dist/
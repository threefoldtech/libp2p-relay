#!/bin/sh

mkdir -p dist
GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.Version=$(git describe)'" -o dist/
cd dist && tar -czf libp2p-relay_$(git describe | sed 's/^v//')_linux_amd64.tar.gz libp2p-relay

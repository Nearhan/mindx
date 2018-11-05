#!/bin/bash

set -e

# build proto
protoc --go_out=plugins=grpc:. proto/*.proto

# build binaries & docker
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o app cmd/proto/generator/main.go 
docker build -t generator -f Dockerfile .
rm app

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o app cmd/proto/processor/main.go 
docker build -t processor -f Dockerfile .
rm app
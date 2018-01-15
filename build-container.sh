#!/bin/sh
docker run --rm -v $GOPATH:/go -w /go/src/psychic-rat golang:latest go build -a -ldflags '-extldflags "-static"' psychic-rat/cmd/server
cp -r server pr.dat res/ docker/
docker build -t pr docker

#!/bin/sh
set -e
docker run --rm -v $GOPATH:/go -w /go/src/psychic-rat golang:1.9.2 go build -a -ldflags '-extldflags "-static"' -o server.linux psychic-rat/cmd/server 
rm -rf docker/res
cp -r server.linux pr.dat res docker/
[ -d content ] && cp content/*tmpl docker/res/tmpl/
docker build -t pr docker

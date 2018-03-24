#!/bin/sh
docker run --rm -v $GOPATH:/go -e GOOS=linux -e CGO_ENABLED=0 -w /go/src/psychic-rat golang:1.9.2-alpine go build -o server.linux psychic-rat/cmd/server 
rm -rf docker/res
cp -r server.linux pr.dat res docker/
[ -d content ] && cp content/*tmpl docker/res/tmpl/
docker build -t pr docker

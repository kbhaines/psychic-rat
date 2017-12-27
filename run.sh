#!/bin/sh
set -e

[ -d rundir ] || mkdir -p rundir/res
cp res/*tmpl rundir/res/
[ -f pr.dat ] || go run cmd/gendb/main.go > /dev/null
cp pr.dat rundir
cd rundir
go build psychic-rat/cmd/server && ./server $*

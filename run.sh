#!/bin/sh
set -e -x

[ -d rundir ] || mkdir -p rundir/res
cp res/*tmpl rundir/res/
cp -r res/css rundir/css
cp -r res/images rundir/images
cp -r res/js rundir/js
[ -f pr.dat ] || go run cmd/gendb/main.go > /dev/null
cp pr.dat rundir
cd rundir
go build psychic-rat/cmd/server && ./server $*

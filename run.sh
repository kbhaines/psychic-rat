#!/bin/sh
set -e -x

[ -d rundir ] || mkdir -p rundir/res
rm -rf rundir/res
cp -r res rundir/res/
ln -sf `pwd`/content/*.tmpl rundir/res/tmpl/
[ -f pr.dat ] || go run cmd/gendb/main.go > /dev/null
cp pr.dat rundir
cd rundir
go build psychic-rat/cmd/server && ./server $*

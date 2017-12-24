#!/bin/sh
set -e

[ -d rundir ] || mkdir -p rundir/res
cp web/res/*tmpl rundir/res/
cp pr.dat rundir
cd rundir
go build psychic-rat/cmd/server && ./server

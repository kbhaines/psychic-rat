#!/bin/sh

set -e
set -x

HOST=$1
DATE=`date +"%y%m%d-%H%M"`

ssh $HOST docker kill --signal SIGUSR1 pr
ssh $HOST docker exec pr cat /pr.dat.bak | gzip > db-$DATE.dat.gz


#!/bin/sh

cp web/res/*tmpl .
go build psychic-rat/cmd/server && ./server

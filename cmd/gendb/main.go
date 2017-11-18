package main

import (
	"flag"
	"psychic-rat/datagen"
	"psychic-rat/sqldb"
)

func main() {

	size := flag.Int("size", 10000, "db overall size, in records")
	file := flag.String("file", "pr.dat", "name of DB file to output")
	flag.Parse()

	db, err := sqldb.NewDB(*file)
	if err != nil {
		panic("could not initialise db:" + err.Error())
	}

	datagen.Generate(db, *size)
}

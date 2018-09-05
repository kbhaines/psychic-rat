package main

import (
	"flag"
	"os"
	"psychic-rat/datagen"
	"psychic-rat/sqldb"
	"runtime/pprof"
)

func main() {
	pprof.StartCPUProfile(os.Stdout)

	size := flag.Int("size", 10000, "db overall size, in records")
	file := flag.String("file", "pr.dat", "name of DB file to output")
	flag.Parse()

	sql, err := sqldb.NewSqliteDB(*file)
	db := sqldb.NewDB(sql)
	if err != nil {
		panic("could not initialise db:" + err.Error())
	}

	datagen.Generate(db, *size)
	db.Close()

	pprof.StopCPUProfile()
}

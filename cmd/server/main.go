package main

import (
	"flag"
	"net/http"
	"psychic-rat/auth0"
	"psychic-rat/sqldb"
	"psychic-rat/web"
	"psychic-rat/web/admin"
	"psychic-rat/web/tmpl"

	"github.com/gorilla/context"
)

var (
	flags struct {
		enableAuth0, sqldb bool
	}
)

func main() {
	flag.BoolVar(&flags.enableAuth0, "auth0", false, "enable auth0 function")
	flag.BoolVar(&flags.sqldb, "sqldb", false, "enable real database")
	flag.Parse()

	http.ListenAndServe("localhost:8080", context.ClearHandler(web.Handler()))
}

func init() {
	db, err := sqldb.OpenDB("pr.dat")
	if err != nil {
		panic("unable to init db: " + err.Error())
	}

	apis := web.APIS{
		Company: db,
		Item:    db,
		NewItem: db,
		Pledge:  db,
		User:    db,
	}
	web.Init(apis)
	auth0.Init(apis.User)
	admin.Init(db, db, db, db)
	tmpl.Init("res/")
}

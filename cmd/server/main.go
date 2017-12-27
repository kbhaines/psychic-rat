package main

import (
	"flag"
	"net/http"
	"psychic-rat/auth0"
	"psychic-rat/authsimple"
	"psychic-rat/sqldb"
	"psychic-rat/web"
	"psychic-rat/web/admin"
	"psychic-rat/web/pub"
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

	apis := pub.APIS{
		Item:    db,
		NewItem: db,
		Pledge:  db,
		User:    db,
	}
	auth0.Init(db)
	renderer := tmpl.NewRenderer("res/")
	pub.Init(apis, flags.enableAuth0, authsimple.NewAuthSimple(db, renderer), renderer)
	admin.Init(db, db, db, db, renderer)
}

package main

import (
	"flag"
	"log"
	"net/http"
	"os"
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

	initModules()
	err := http.ListenAndServe("localhost:8080", context.ClearHandler(web.Handler()))
	if err != nil {
		log.Fatalf("web server aborted: %v", err)
	}
}

func initModules() {
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
	var authHandler pub.AuthHandler
	if flags.enableAuth0 {
		authHandler = auth0.NewAuth0Handler(renderer, os.Getenv("AUTH0_DOMAIN"), os.Getenv("AUTH0_CALLBACK_URL"), os.Getenv("AUTH0_CLIENT_ID"))
	} else {
		authHandler = authsimple.NewAuthSimple(db, renderer)
	}
	pub.Init(apis, authHandler, renderer)
	admin.Init(db, db, db, db, authHandler, renderer)
}

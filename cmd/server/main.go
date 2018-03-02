package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"psychic-rat/auth"
	"psychic-rat/auth0"
	"psychic-rat/authsimple"
	"psychic-rat/sess"
	"psychic-rat/sqldb"
	"psychic-rat/web"
	"psychic-rat/web/admin"
	"psychic-rat/web/pub"
	"psychic-rat/web/tmpl"

	"github.com/gorilla/context"
)

var (
	flags struct {
		enableAuth0    bool
		sqldb          bool
		cacheTemplates bool
		listenOn       string
		mockCSRF       bool
	}
)

func main() {
	flag.StringVar(&flags.listenOn, "listen", "localhost:8080", "interface:port to listen on")
	flag.BoolVar(&flags.enableAuth0, "auth0", false, "enable auth0 function")
	flag.BoolVar(&flags.sqldb, "sqldb", false, "enable real database")
	flag.BoolVar(&flags.cacheTemplates, "cache-templates", false, "enable template caching")
	flag.BoolVar(&flags.mockCSRF, "mockcsrf", false, "mock CSRF token")
	flag.Parse()

	initModules()
	err := http.ListenAndServe(flags.listenOn, context.ClearHandler(web.Handler()))
	if err != nil {
		log.Fatalf("web server aborted: %v", err)
	}
}

func initModules() {
	db, err := sqldb.OpenDB("pr.dat")
	if err != nil {
		panic("unable to init db: " + err.Error())
	}

	sess.Init(flags.mockCSRF)
	renderer := tmpl.NewRenderer("res/tmpl", flags.cacheTemplates)
	var authHandler pub.AuthHandler
	if flags.enableAuth0 {
		authHandler = auth0.NewAuth0Handler(renderer, "me", "http://localhost:8080/auth/facebook", os.Getenv("FACEBOOK_CLIENT_ID"))
	} else {
		authHandler = authsimple.NewAuthSimple(db, renderer)
	}
	serverURL := os.Getenv("SERVER_URL")
	if serverURL == "" {
		serverURL = "http://localhost:8080/"
	}
	auth.Init(db, serverURL+"callback")
	//TODO: take authHandler out...
	pub.Init(db, db, db, authHandler, renderer)
	admin.Init(db, db, db, db, authHandler, renderer)
}

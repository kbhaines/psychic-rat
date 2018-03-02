package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"psychic-rat/auth"
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
		sqldb          bool
		cacheTemplates bool
		listenOn       string
		mockCSRF       bool
	}
)

func main() {
	flag.StringVar(&flags.listenOn, "listen", "localhost:8080", "interface:port to listen on")
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
	serverURL := os.Getenv("SERVER_URL")
	if serverURL == "" {
		serverURL = "http://localhost:8080/"
	}
	auth.Init(db, serverURL+"callback")

	userHandler := auth.NewUserHandler()
	renderer := tmpl.NewRenderer("res/tmpl", flags.cacheTemplates)
	pub.Init(db, db, db, userHandler, renderer)
	admin.Init(db, db, db, db, userHandler, renderer)
}

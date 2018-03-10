package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"psychic-rat/auth"
	"psychic-rat/auth/basic"
	"psychic-rat/auth/facebook"
	"psychic-rat/auth/gplus"
	"psychic-rat/auth/recaptcha"
	"psychic-rat/auth/twitter"
	"psychic-rat/limit"
	"psychic-rat/sess"
	"psychic-rat/sqldb"
	"psychic-rat/types"
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
		basicAuth      bool
	}
)

type UserHandler interface {
	GetLoggedInUser(*http.Request) (*types.User, error)
	GetUserCSRF(http.ResponseWriter, *http.Request) (string, error)
	LogOut(http.ResponseWriter, *http.Request) error
	VerifyUserCSRF(*http.Request, string) error
}

func main() {
	flag.StringVar(&flags.listenOn, "listen", "localhost:8080", "interface:port to listen on")
	flag.BoolVar(&flags.sqldb, "sqldb", false, "enable real database")
	flag.BoolVar(&flags.cacheTemplates, "cache-templates", false, "enable template caching")
	flag.BoolVar(&flags.basicAuth, "basicauth", false, "enable basic auth mode for testing")
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

	sess.Init(flags.basicAuth)
	serverURL := os.Getenv("SERVER_URL")
	if serverURL == "" {
		serverURL = "http://localhost:8080/"
		log.Printf("WARNING: serving from %s", serverURL)
	}

	callbackURL := serverURL + "callback"

	var (
		userHandler   UserHandler
		authProviders map[string]auth.AuthHandler
	)

	if !flags.basicAuth {
		userHandler = auth.NewUserHandler()
		authProviders = map[string]auth.AuthHandler{
			"facebook": facebook.New(callbackURL+"?p=facebook", userHandler),
			"twitter":  twitter.New(callbackURL + "?p=twitter"),
			"gplus":    gplus.New(callbackURL+"?p=gplus", userHandler),
		}
	} else {
		authProviders = map[string]auth.AuthHandler{
			"basic": basic.New(callbackURL + "?p=basic"),
		}
		userHandler = basic.NewUserHandler()
		log.Printf("WARNING: using basic auth")
	}

	auth.Init(db, authProviders)
	web.Init(userHandler, limit.New())
	renderer := tmpl.NewRenderer("res/tmpl", flags.cacheTemplates)
	pub.Init(db, db, db, userHandler, renderer, recaptcha.New(os.Getenv("RECAPTCHA_SECRET")))
	admin.Init(db, db, db, db, userHandler, renderer)
}

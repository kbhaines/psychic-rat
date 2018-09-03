package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
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
	"strings"
	"syscall"

	ctxt "context"

	_ "net/http/pprof"

	"github.com/gorilla/context"
)

var (
	flags struct {
		sqldb          bool
		cacheTemplates bool
		listenOn       string
		listenSSL      string
		basicAuth      bool
		limit          string
		certFile       string
		keyFile        string
		memprofile     bool
	}
)

type (
	UserHandler interface {
		GetLoggedInUser(*http.Request) (*types.User, error)
		GetUserCSRF(http.ResponseWriter, *http.Request) (string, error)
		LogOut(http.ResponseWriter, *http.Request) error
		VerifyUserCSRF(*http.Request, string) error
	}

	fakeCaptcha struct{}
)

func main() {
	flag.StringVar(&flags.listenOn, "listen", "localhost:8080", "interface:port to listen on")
	flag.StringVar(&flags.listenSSL, "listenSSL", "localhost:4444", "interface:port to listen on for SSL")
	flag.BoolVar(&flags.sqldb, "sqldb", false, "enable real database")
	flag.BoolVar(&flags.cacheTemplates, "cache-templates", false, "enable template caching")
	flag.BoolVar(&flags.basicAuth, "basicauth", false, "enable basic auth mode for testing")
	flag.StringVar(&flags.limit, "limit", "30,10,5", "rate-limit bucket specification")
	flag.StringVar(&flags.certFile, "cert", "/cert/fullchain.pem", "SSL certificate filename")
	flag.StringVar(&flags.keyFile, "key", "/cert/privkey.pem", "SSL key filename")
	flag.BoolVar(&flags.memprofile, "profile", false, "serve profiling")
	flag.Parse()

	if flags.memprofile {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	db := initModules()
	shutdownChan := initSignalHandler()
	initBackupHandler("pr.dat", "pr.dat.bak")
	runServer(shutdownChan)
	db.Close()
}

func initModules() *sqldb.DB {
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
		authProviders map[string]auth.AuthHandler
		humanTest     pub.HumanTester
	)

	userHandler := auth.NewUserHandler()
	if !flags.basicAuth {
		authProviders = map[string]auth.AuthHandler{
			"facebook": facebook.New(callbackURL+"?p=facebook", userHandler),
			"twitter":  twitter.New(callbackURL + "?p=twitter"),
			"gplus":    gplus.New(callbackURL+"?p=gplus", userHandler),
		}
		humanTest = recaptcha.New(os.Getenv("RECAPTCHA_SECRET"))
	} else {
		authProviders = map[string]auth.AuthHandler{
			"basic": basic.New(callbackURL + "?p=basic"),
		}
		humanTest = &fakeCaptcha{}
		log.Printf("WARNING: using basic auth, disabled captcha")
	}

	auth.Init(db, authProviders)

	var max, increment, interval int
	_, err = fmt.Sscanf(flags.limit, "%d,%d,%d", &max, &increment, &interval)
	if err != nil {
		panic(err)
	}

	web.Init(userHandler, limit.New(max, increment, interval, idGenerator))
	renderer := tmpl.NewRenderer("res/tmpl", flags.cacheTemplates)
	pub.Init(db, db, db, userHandler, renderer, humanTest)
	admin.Init(db, db, db, db, userHandler, renderer)
	return db
}

func initSignalHandler() chan bool {
	sigs := make(chan os.Signal, 1)
	shutdown := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-sigs
		log.Printf("received signal %v", s)
		shutdown <- true
	}()

	return shutdown
}

func initBackupHandler(dbName string, backupName string) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1)
	go func() {
		for {
			s := <-sigs
			log.Printf("received backup request %v", s)
			err := sqldb.Backup(dbName, backupName)
			if err != nil {
				log.Printf("backup error: %v", err)
			}
			log.Printf("backup completed")
		}
	}()
}

func runServer(shutdown chan bool) {
	srvSSL := &http.Server{Addr: flags.listenSSL, Handler: context.ClearHandler(web.Handler())}
	go func() {
		err := srvSSL.ListenAndServeTLS(flags.certFile, flags.keyFile)
		log.Printf("SSL web server shutdown: %v", err)
	}()
	srv := &http.Server{Addr: flags.listenOn, Handler: context.ClearHandler(web.Handler())}
	go func() {
		err := srv.ListenAndServe()
		log.Printf("web server shutdown: %v", err)
	}()

	<-shutdown
	srv.Shutdown(ctxt.Background())
	srvSSL.Shutdown(ctxt.Background())
}

func idGenerator(r *http.Request) string { return strings.Split(r.RemoteAddr, ":")[0] }

func (_ *fakeCaptcha) IsHuman(url.Values) error { return nil }

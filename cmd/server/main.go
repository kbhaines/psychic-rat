package main

import (
	"flag"
	"net/http"
	"psychic-rat/web"

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

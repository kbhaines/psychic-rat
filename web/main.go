package web

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"psychic-rat/auth0"
	"psychic-rat/log"
	"psychic-rat/sess"
	"psychic-rat/types"
	"psychic-rat/web/admin"
	"psychic-rat/web/dispatch"
	"psychic-rat/web/pub"
)

const (
	HomePage      = "/"
	SignInPage    = "/signin"
	SignOutPage   = "/signout"
	PledgePage    = "/pledge"
	ThanksPage    = "/thanks"
	NewItem       = "/newitem"
	AdminNewItems = "/admin/newitems"
	Callback      = "/callback"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

var (
	uriHandlers = []dispatch.URIHandler{
		{HomePage, pub.HomePageHandler},
		{SignInPage, pub.SignInPageHandler},
		{SignOutPage, pub.SignOutPageHandler},
		{PledgePage, pub.PledgePageHandler},
		{NewItem, pub.NewItemHandler},
		{ThanksPage, pub.ThanksPageHandler},
		{Callback, auth0.CallbackHandler},
		{AdminNewItems, admin.AdminItemHandler},
	}

	flags struct {
		enableAuth0, sqldb bool
	}
)

func Handler() http.Handler {
	hmux := http.NewServeMux()
	handlerForDirs(hmux, "css", "js", "images")
	for _, h := range uriHandlers {
		hmux.HandleFunc(h.URI, addContextValues(logRequest(h.Handler)))
	}
	return hmux
}

func addContextValues(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := rand.Int63()
		user, _ := sess.NewSessionStore(r).Get()
		if user == nil {
			user = &types.User{ID: "<none>"}
		}
		ctx := context.WithValue(r.Context(), "rid", requestID)
		ctx = context.WithValue(ctx, "uid", user.ID)
		r = r.WithContext(ctx)
		next(w, r)
	}
}

func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := fmt.Sprintf("source: %s request: %s %s", r.RemoteAddr, r.Method, r.RequestURI)
		log.Logf(r, request)
		scw := &statusRecorder{w, 200}
		next(scw, r)
		log.Logf(r, "response: %d", scw.status)
	}
}

func handlerForDirs(mux *http.ServeMux, dir ...string) {
	for _, d := range dir {
		mux.Handle("/"+d+"/", http.StripPrefix("/"+d, http.FileServer(http.Dir("res/"+d+"/"))))
	}
}

func (s *statusRecorder) WriteHeader(c int) {
	s.status = c
	s.ResponseWriter.WriteHeader(c)
}

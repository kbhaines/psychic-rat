package web

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"psychic-rat/auth/facebook"
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
		{"/auth/facebook", facebook.CallbackHandler},
	}

	flags struct {
		enableAuth0, sqldb bool
	}
)

func Handler() http.Handler {
	hmux := http.NewServeMux()
	handlerForDirs(hmux, "css", "js", "images")
	for _, h := range uriHandlers {
		hmux.HandleFunc(h.URI, addContextValues(logRequest(csrfProtect(h.Handler))))
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
		request := fmt.Sprintf("request: %s %s source: %s", r.Method, r.RequestURI, r.RemoteAddr)
		log.Logf(r.Context(), request)
		scw := &statusRecorder{w, 200}
		next(scw, r)
		log.Logf(r.Context(), "response: %d", scw.status)
	}
}

func csrfProtect(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			next(w, r)
			return
		}
		err := r.ParseForm()
		if err != nil {
			log.Errorf(r.Context(), "csrfProtect: failed to parse form: %v", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		err = sess.NewSessionStore(r).VerifyCSRF(r.FormValue("csrf"))
		if err != nil {
			log.Errorf(r.Context(), "csrfProtect: CSRF failed validation: %v", err)
			http.Error(w, "", http.StatusForbidden)
			return
		}
		log.Logf(r.Context(), "csrf check ok")
		next(w, r)
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

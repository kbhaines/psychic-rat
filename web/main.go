package web

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"psychic-rat/auth"
	"psychic-rat/log"
	"psychic-rat/types"
	"psychic-rat/web/admin"
	"psychic-rat/web/dispatch"
	"psychic-rat/web/pub"

	gorcon "github.com/gorilla/context"
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
	AuthInit      = "/auth"
)

type (
	statusRecorder struct {
		http.ResponseWriter
		status int
	}

	UserHandler interface {
		GetLoggedInUser(*http.Request) (*types.User, error)
		VerifyUserCSRF(*http.Request, string) error
	}

	RateLimiter interface {
		CheckLimit(*http.Request) error
	}
)

var (
	uriHandlers = []dispatch.URIHandler{
		{HomePage, pub.HomePageHandler},
		{SignInPage, pub.SignInPageHandler},
		{SignOutPage, pub.SignOutPageHandler},
		{PledgePage, pub.PledgePageHandler},
		{NewItem, pub.NewItemHandler},
		{ThanksPage, pub.ThanksPageHandler},
		{AdminNewItems, admin.AdminItemHandler},
		{AuthInit, auth.AuthInit},
		{Callback, auth.CallbackHandler},
	}

	userHandler UserHandler
	rateLimiter RateLimiter

	flags struct {
		enableAuth0, sqldb bool
	}
)

func Init(u UserHandler, r RateLimiter) {
	userHandler = u
	rateLimiter = r
}

func Handler() http.Handler {
	hmux := http.NewServeMux()
	handlerForDirs(hmux, "css", "js", "images", ".well-known")
	for _, h := range uriHandlers {
		hmux.HandleFunc(h.URI, addContextValues(logRequest(rateLimit(csrfProtect(h.Handler)))))
	}
	return hmux
}

func handlerForDirs(mux *http.ServeMux, dir ...string) {
	for _, d := range dir {
		baseHandler := http.StripPrefix("/"+d, http.FileServer(http.Dir("res/"+d+"/")))
		mux.HandleFunc("/"+d+"/", addContextValues(logRequest(rateLimit(baseHandler.ServeHTTP))))
	}
}

func addContextValues(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := rand.Int63()
		user, err := userHandler.GetLoggedInUser(r)
		if user == nil || err != nil {
			user = &types.User{ID: "<none>"}
		}
		ctx := context.WithValue(r.Context(), "rid", requestID)
		ctx = context.WithValue(ctx, "uid", user.ID)
		r = r.WithContext(ctx)
		next(w, r)

		// Clear the request we created here - otherwise there's a leak as Gorilla
		// keeps the request pointer in a map. The original request is already
		// handled by the ClearHandler call from main.runServer
		gorcon.Clear(r)
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

func rateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := rateLimiter.CheckLimit(r); err != nil {
			log.Errorf(r.Context(), "rate limited: %v", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "<html><meta http-equiv=refresh content=5;url=%s /></html>server has hit resource limits; retrying automatically in 5 seconds",
				r.RequestURI)
			return
		}
		next(w, r)
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
		err = userHandler.VerifyUserCSRF(r, r.FormValue("csrf"))
		if err != nil {
			log.Errorf(r.Context(), "csrfProtect: CSRF failed validation: %v", err)
			http.Error(w, "", http.StatusForbidden)
			return
		}
		log.Logf(r.Context(), "csrf check ok")
		next(w, r)
	}
}

func (s *statusRecorder) WriteHeader(c int) {
	s.status = c
	s.ResponseWriter.WriteHeader(c)
}

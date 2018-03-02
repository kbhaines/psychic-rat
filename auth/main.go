package auth

import (
	"net/http"
	"psychic-rat/auth/facebook"
	"psychic-rat/auth/gplus"
	"psychic-rat/auth/twitter"
	"psychic-rat/log"
	"psychic-rat/sess"
	"psychic-rat/types"
)

type (
	UserAPI interface {
		GetUser(id string) (*types.User, error)
		AddUser(types.User) error
	}

	AuthHandler interface {
		BeginAuth(w http.ResponseWriter, r *http.Request) (string, error)
		Callback(w http.ResponseWriter, r *http.Request) (*types.User, error)
	}
)

var (
	authProviders map[string]AuthHandler
	userAPI       UserAPI
)

func Init(u UserAPI, callbackURL string) {
	userAPI = u

	//TODO: inject this dependency
	authProviders = map[string]AuthHandler{
		"facebook": facebook.New(callbackURL + "?p=facebook"),
		"twitter":  twitter.New(callbackURL + "?p=twitter"),
		"gplus":    gplus.New(callbackURL + "?p=gplus"),
	}
}

func AuthInit(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("p")
	handler, ok := authProviders[provider]
	if !ok {
		log.Errorf(r.Context(), "provider %s not found", provider)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	redirect, err := handler.BeginAuth(w, r)
	if err != nil {
		log.Errorf(r.Context(), "unable to start auth process for %s: %v", provider, err)
		http.Error(w, "could not start auth", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("p")
	handler, ok := authProviders[provider]
	if !ok {
		log.Errorf(r.Context(), "callback provider %s not found", provider)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := handler.Callback(w, r)
	if err != nil {
		log.Errorf(r.Context(), "could not handle callback: %v", err)
		http.Error(w, "auth error has occurred", http.StatusInternalServerError)
		return
	}

	user, err = addUserIfNotExists(user)
	if err != nil {
		log.Errorf(r.Context(), "unable to create a user %v :%v", user, err)
		http.Error(w, "auth error has occurred", http.StatusInternalServerError)
		return
	}

	err = sess.NewSessionStore(r).Save(user, w)
	if err != nil {
		log.Errorf(r.Context(), "unable to save user into session: %v", err)
		http.Error(w, "auth error has occurred", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func addUserIfNotExists(u *types.User) (*types.User, error) {
	existing, err := userAPI.GetUser(u.ID)
	if err != nil {
		return u, userAPI.AddUser(*u)
	}
	return existing, nil
}

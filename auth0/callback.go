package auth0

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"psychic-rat/sess"
	"psychic-rat/types"

	"golang.org/x/oauth2"
)

type UserAPI interface {
	GetUser(id string) (*types.User, error)
	AddUser(types.User) error
}

var (
	userAPI UserAPI
)

func Init(a UserAPI) {
	userAPI = a
}

// TODO: taken from Auth2's sample, refactor

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	domain := os.Getenv("AUTH0_DOMAIN")
	conf := &oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Scopes:       []string{"openid", "profile", "user_metadata"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + domain + "/authorize",
			TokenURL: "https://" + domain + "/oauth/token",
		},
	}

	code := r.URL.Query().Get("code")

	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("err = %+v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := conf.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://" + domain + "/userinfo")
	if err != nil {
		log.Printf("err = %+v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	raw, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Printf("err = %+v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	if err = json.Unmarshal(raw, &profile); err != nil {
		log.Printf("err = %+v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("profile = %+v\n", profile)

	//gob.Register(map[string]interface{}{})
	//session.Values["id_token"] = token.Extra("id_token")
	//session.Values["access_token"] = token.AccessToken
	//session.Values["profile"] = profile
	//session.Values["userId"] = userId

	userId := profile["sub"].(string)
	userRecord, error := userAPI.GetUser(userId)
	if error != nil {
		userRecord = &types.User{
			ID:        userId,
			Fullname:  profile["name"].(string),
			FirstName: profile["given_name"].(string),
			Country:   profile["locale"].(string),
		}
		err := userAPI.AddUser(*userRecord)
		if err != nil {
			log.Fatal("unable to create a user %v :%v", userRecord, err)
			return
		}

	}
	store := sess.NewSessionStore(r, w)
	store.Save(*userRecord)
	if err != nil {
		log.Fatal("unable to save user into session: %v", err)
		return
	}

	// Redirect to logged in page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type (
	Auth0Handler struct {
		renderer    Renderer
		domain      string
		callbackURL string
		clientID    string
	}

	Renderer interface {
		Render(http.ResponseWriter, string, interface{}) error
	}
)

func NewAuth0Handler(rendr Renderer, dom string, callback string, client string) *Auth0Handler {
	return &Auth0Handler{
		renderer:    rendr,
		domain:      dom,
		callbackURL: callback,
		clientID:    client,
	}
}

func (a *Auth0Handler) Handler(writer http.ResponseWriter, request *http.Request) {
	vars := struct {
		Auth0Domain      string
		Auth0CallbackURL string
		Auth0ClientId    string
		User             types.User
	}{
		Auth0Domain:      a.domain,
		Auth0CallbackURL: a.callbackURL,
		Auth0ClientId:    a.clientID,
	}
	log.Printf("vars = %+v\n", vars)
	a.renderer.Render(writer, "signin-auth0.html.tmpl", vars)
}

func (a *Auth0Handler) GetLoggedInUser(r *http.Request) (*types.User, error) {
	// TODO: nil is a smell. StoreReader/Writer interfaces.
	s := sess.NewSessionStore(r, nil)
	user, err := s.Get()
	if err != nil {
		return nil, err
	}
	return user, nil
}

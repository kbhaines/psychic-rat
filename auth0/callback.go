package auth0

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"psychic-rat/api/rest"
	"psychic-rat/mdl"
	"psychic-rat/sess"

	"golang.org/x/oauth2"
)

type UserAPI interface {
	GetById(mdl.Id) (*mdl.UserRecord, error)
	Create(mdl.UserRecord) error
}

var (
	userAPI   UserAPI
	serverURL string
)

func Init(a UserAPI, server string) {
	userAPI = a
	serverURL = server
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Scopes:       []string{"openid", "profile", "user_metadata"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  serverURL + "/authorize",
			TokenURL: serverURL + "/oauth/token",
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
	resp, err := client.Get(serverURL + "/userinfo")
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

	userId := mdl.Id(profile["sub"].(string))
	userRecord, error := userAPI.GetById(userId)
	if error != nil {
		userRecord = &mdl.UserRecord{
			Id:        userId,
			Fullname:  profile["name"].(string),
			FirstName: profile["given_name"].(string),
			Country:   profile["locale"].(string),
		}
		err := userAPI.Create(*userRecord)
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
	http.Redirect(w, r, rest.HomePage, http.StatusSeeOther)
}

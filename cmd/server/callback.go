package main

import (
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"psychic-rat/api/rest"
	"psychic-rat/mdl"

	"golang.org/x/oauth2"
)

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Getting now the userInfo
	client := conf.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://" + domain + "/userinfo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	raw, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	if err = json.Unmarshal(raw, &profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("profile = %+v\n", profile)

	gob.Register(map[string]interface{}{})
	gob.Register(mdl.Id(""))
	gob.Register(mdl.UserRecord{})

	session, err := auth0Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["id_token"] = token.Extra("id_token")
	session.Values["access_token"] = token.AccessToken
	session.Values["profile"] = profile
	userId := mdl.Id(profile["sub"].(string))
	session.Values["userId"] = userId

	userRecord, err := apis.User.GetById(userId)
	if err != nil {
		userRecord = &mdl.UserRecord{
			Id:        userId,
			Fullname:  profile["name"].(string),
			FirstName: profile["given_name"].(string),
			Country:   profile["locale"].(string),
		}
		apis.User.Create(*userRecord)
	}
	session.Values["userRecord"] = *userRecord

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to logged in page
	http.Redirect(w, r, rest.HomePage, http.StatusSeeOther)

}
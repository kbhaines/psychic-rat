package sess

import (
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"psychic-rat/types"
	"strings"

	"github.com/gorilla/sessions"
)

type (
	SessionStore struct {
		r     *http.Request
		w     http.ResponseWriter
		store sessions.Store
	}
)

const sessionVar = "auth"

var cookieKeys [][]byte

func init() {
	gob.Register(types.User{})
	cookieKeys = getKeys()
}

func NewSessionStore(r *http.Request, w http.ResponseWriter) *SessionStore {
	return &SessionStore{r: r, w: w, store: sessions.NewCookieStore(cookieKeys...)}
}

func (s *SessionStore) Get() (*types.User, error) {
	session, err := s.store.Get(s.r, sessionVar)
	if err != nil {
		log.Printf("Get: error retrieving, trying re-write: %v", err)
		return nil, nil
	}
	userFromSession, found := session.Values["userRecord"]
	if !found {
		return nil, nil
	}
	userRecord, ok := userFromSession.(types.User)
	if !ok {
		return nil, fmt.Errorf("Get: conversion error: %v", err)
	}
	log.Printf("loaded user %v from session", userRecord)
	return &userRecord, nil
}

func (s *SessionStore) Save(user *types.User) error {
	session, err := s.store.Get(s.r, sessionVar)
	if err != nil {
		log.Printf("save: cannot retrieve from store, rewriting: %v", err)
		s.store.Save(s.r, s.w, session)
	}
	if user != nil {
		session.Values["userRecord"] = *user
	} else {
		delete(session.Values, "userRecord")
	}

	if err := session.Save(s.r, s.w); err != nil {
		http.Error(s.w, err.Error(), http.StatusInternalServerError)
		return err
	}
	log.Printf("saved user %v in session", user)
	return nil
}

func (s *SessionStore) Request() *http.Request      { return s.r }
func (s *SessionStore) Writer() http.ResponseWriter { return s.w }

func getKeys() [][]byte {
	keys := os.Getenv("COOKIE_KEYS")
	if keys == "" {
		log.Printf("Warning: using default cookie keys")
		keys = "defaultnotsafe"
	}
	results := make([][]byte, 0, len(keys)*2)
	for _, key := range strings.Split(keys, ",") {
		keySha := sha256.Sum256([]byte(key))
		log.Printf("key,keySha = %+v %+v\n", key, keySha)
		results = append(results, []byte(key), keySha[:])
	}
	return results
}

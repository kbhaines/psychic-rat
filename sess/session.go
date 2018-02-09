package sess

import (
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	syslog "log"
	"net/http"
	"os"
	"psychic-rat/log"
	"psychic-rat/types"
	"strings"

	"github.com/gorilla/sessions"
)

type (
	SessionStore struct {
		r     *http.Request
		store sessions.Store
	}
)

const sessionVar = "auth"

var cookieKeys [][]byte

func init() {
	gob.Register(types.User{})
	cookieKeys = getKeys()
}

func NewSessionStore(r *http.Request) *SessionStore {
	return &SessionStore{r: r, store: sessions.NewCookieStore(cookieKeys...)}
}

func (s *SessionStore) Get() (*types.User, error) {
	session, err := s.store.Get(s.r, sessionVar)
	if err != nil {
		log.Logf(s.r, "Get: error retrieving: %v", err)
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
	return &userRecord, nil
}

func (s *SessionStore) Save(user *types.User, w http.ResponseWriter) error {
	session, err := s.store.Get(s.r, sessionVar)
	if err != nil {
		log.Logf(s.r, "save: cannot retrieve from store, rewriting: %v", err)
		s.store.Save(s.r, w, session)
	}
	if user != nil {
		session.Values["userRecord"] = *user
	} else {
		delete(session.Values, "userRecord")
	}

	if err := session.Save(s.r, w); err != nil {
		return err
	}
	log.Logf(s.r, "saved user %v in session", user)
	return nil
}

func (s *SessionStore) Request() *http.Request { return s.r }

func getKeys() [][]byte {
	keys := os.Getenv("COOKIE_KEYS")
	if keys == "" {
		syslog.Printf("Warning: using default cookie keys")
		keys = "defaultnotsafe"
	}
	results := make([][]byte, 0, len(keys)*2)
	for _, key := range strings.Split(keys, ",") {
		keySha := sha256.Sum256([]byte(key))
		syslog.Printf("key,keySha = %+v %+v\n", key, keySha)
		results = append(results, []byte(key), keySha[:])
	}
	return results
}

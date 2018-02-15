package sess

import (
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	syslog "log"
	"math/rand"
	"net/http"
	"os"
	"psychic-rat/log"
	"psychic-rat/types"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/sessions"
)

type (
	SessionStore struct {
		r     *http.Request
		store sessions.Store
	}
)

const (
	sessionVar = "auth"
	userVar    = "userRecord"
	csrfVar    = "csrf"
)

var (
	cookieKeys    [][]byte
	csrfGenerator = func() string { return strconv.FormatInt(rand.Int63(), 36) }
)

func init() {
	gob.Register(types.User{})
	cookieKeys = getKeys()
	rand.Seed(time.Now().UTC().UnixNano())
}

func Init(useMockCSRF bool) {
	csrfGenerator = func() string { return "mockcsrf" }
}

func NewSessionStore(r *http.Request) *SessionStore {
	return &SessionStore{r: r, store: sessions.NewCookieStore(cookieKeys...)}
}

func (s *SessionStore) Get() (*types.User, error) {
	session, err := s.store.Get(s.r, sessionVar)
	if err != nil {
		log.Logf(s.r.Context(), "Get: error retrieving: %v", err)
		return nil, nil
	}
	userFromSession, found := session.Values[userVar]
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
		log.Logf(s.r.Context(), "save: cannot retrieve from store, rewriting: %v", err)
		s.store.Save(s.r, w, session)
	}
	if user != nil {
		session.Values[userVar] = *user
	} else {
		delete(session.Values, userVar)
	}

	if err := session.Save(s.r, w); err != nil {
		return err
	}
	log.Logf(s.r.Context(), "saved user %v in session", user)
	return nil
}

func (s *SessionStore) SetCSRF(w http.ResponseWriter) (string, error) {
	session, err := s.store.Get(s.r, sessionVar)
	if err != nil {
		log.Logf(s.r.Context(), "SetCSRF: cannot retrieve from store, rewriting: %v", err)
		s.store.Save(s.r, w, session)
	}
	token := csrfGenerator()
	session.Values[csrfVar] = token
	if err := session.Save(s.r, w); err != nil {
		return "", fmt.Errorf("SetCSRF: could not set csrf: %v", err)
	}
	log.Logf(s.r.Context(), "saved csrf in session")
	return token, nil
}

func (s *SessionStore) VerifyCSRF(token string) error {
	session, err := s.store.Get(s.r, sessionVar)
	if err != nil {
		log.Logf(s.r.Context(), "VerifyCSRF: cannot retrieve from store: %v", err)
		return err
	}
	csrfFromSession, found := session.Values[csrfVar]
	if !found {
		return fmt.Errorf("no csrf token in session")
	}
	csrf, ok := csrfFromSession.(string)
	if !ok {
		return fmt.Errorf("VerifyCSRF: conversion error: %v", err)
	}
	if csrf != token {
		return fmt.Errorf("token mismatch, expected %s got %s in POST", csrf, token)
	}
	return nil
}

func (s *SessionStore) Request() *http.Request { return s.r }

func getKeys() [][]byte {
	keys := os.Getenv("COOKIE_KEYS")
	if keys == "" {
		syslog.Printf("WARNING: using default cookie keys")
		keys = "defaultnotsafe"
	}
	results := make([][]byte, 0, len(keys)*2)
	for _, key := range strings.Split(keys, ",") {
		keySha := sha256.Sum256([]byte(key))
		results = append(results, []byte(key), keySha[:])
	}
	return results
}

package sess

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"psychic-rat/mdl"

	"github.com/gorilla/sessions"
)

func init() {
	gob.Register(mdl.ID(0))
	gob.Register(mdl.User{})
}

func NewSessionStore(r *http.Request, w http.ResponseWriter) *SessionStore {
	return &SessionStore{r: r, w: w, store: sessions.NewCookieStore([]byte("something-very-secret"))}
}

type (
	SessionStore struct {
		r     *http.Request
		w     http.ResponseWriter
		store sessions.Store
	}
)

func (s *SessionStore) Get() (*mdl.User, error) {
	session, err := s.store.Get(s.r, "auth-session")
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve from store: %v", err)
	}
	userFromSession, found := session.Values["userRecord"]
	if !found {
		return nil, nil
	}
	userRecord, ok := userFromSession.(mdl.User)
	if !ok {
		return nil, fmt.Errorf("conversion error: %v", err)
	}
	log.Printf("loaded user %v from session", userRecord)
	return &userRecord, nil
}

func (s *SessionStore) Save(user mdl.User) error {
	session, err := s.store.Get(s.r, "auth-session")
	if err != nil {
		return fmt.Errorf("cannot retrieve from store: %v", err)
	}
	session.Values["userRecord"] = user
	if err := session.Save(s.r, s.w); err != nil {
		http.Error(s.w, err.Error(), http.StatusInternalServerError)
		return err
	}
	log.Printf("saved user %v in session", user)
	return nil
}

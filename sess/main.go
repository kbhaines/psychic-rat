package sess

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"psychic-rat/types"

	"github.com/gorilla/sessions"
)

func init() {
	gob.Register(types.User{})
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

func (s *SessionStore) Get() (*types.User, error) {
	session, err := s.store.Get(s.r, "auth-session")
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve from store: %v", err)
	}
	userFromSession, found := session.Values["userRecord"]
	if !found {
		return nil, nil
	}
	userRecord, ok := userFromSession.(types.User)
	if !ok {
		return nil, fmt.Errorf("conversion error: %v", err)
	}
	log.Printf("loaded user %v from session", userRecord)
	return &userRecord, nil
}

func (s *SessionStore) Save(user types.User) error {
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

func (s *SessionStore) Request() *http.Request      { return s.r }
func (s *SessionStore) Writer() http.ResponseWriter { return s.w }

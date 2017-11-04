package sess

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"psychic-rat/mdl"

	"github.com/gorilla/sessions"
)

func NewSessionStore(r *http.Request, w http.ResponseWriter) *SessionStore {
	return &SessionStore{r: r, w: w}
}

var store = sessions.NewCookieStore([]byte("something-very-secret"))

type (
	SessionStore struct {
		r *http.Request
		w http.ResponseWriter
	}
)

func (s *SessionStore) Get() (*mdl.UserRecord, error) {
	session, err := store.Get(s.r, "auth-session")
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve from store: %v", err)
	}
	userFromSession, found := session.Values["userRecord"]
	if !found {
		return nil, nil
	}
	userRecord, ok := userFromSession.(mdl.UserRecord)
	if !ok {
		return nil, fmt.Errorf("conversion error: %v", err)
	}
	log.Println("loaded user %v from session", userRecord)
	return &userRecord, nil
}

func (s *SessionStore) Save(user mdl.UserRecord) error {
	session, err := store.Get(s.r, "auth-session")
	if err != nil {
		return fmt.Errorf("cannot retrieve from store: %v", err)
	}
	session.Values["userRecord"] = user
	if err := session.Save(s.r, s.w); err != nil {
		http.Error(s.w, err.Error(), http.StatusInternalServerError)
		return err
	}
	log.Println("saved user %v in session", user)
	return nil
}

func init() {
	gob.Register(mdl.Id(0))
	gob.Register(mdl.UserRecord{})
}

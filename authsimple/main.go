package authsimple

import (
	"fmt"
	"log"
	"net/http"
	"psychic-rat/sess"
	"psychic-rat/types"
)

type (
	Renderer interface {
		Render(writer http.ResponseWriter, templateName string, variables interface{}) error
	}
	UserAPI interface {
		GetUser(string) (*types.User, error)
	}
	simpleAuthHandler struct {
		userAPI  UserAPI
		renderer Renderer
	}
)

func NewAuthSimple(u UserAPI, r Renderer) *simpleAuthHandler {
	return &simpleAuthHandler{userAPI: u, renderer: r}
}

func (s *simpleAuthHandler) SignIn(sess *sess.SessionStore) {
	if err := s.authUser(sess); err != nil {
		log.Print(err)
		http.Error(sess.Writer(), "authentication failed", http.StatusForbidden)
		return
	}
}

func (s *simpleAuthHandler) Handler(w http.ResponseWriter, r *http.Request) {
	session := sess.NewSessionStore(r, w)
	s.SignIn(session)
	var loggedInUser types.User
	user, err := s.GetLoggedInUser(r)
	if err != nil {
		log.Print(err)
		http.Error(w, "authentication failed", http.StatusForbidden)
		return
	}

	if user != nil {
		loggedInUser = *user
	}

	vars := struct{ User types.User }{loggedInUser}
	s.renderer.Render(w, "signin.html.tmpl", vars)
}

func (_ *simpleAuthHandler) GetLoggedInUser(r *http.Request) (*types.User, error) {
	// TODO: nil is a smell. StoreReader/Writer interfaces.
	s := sess.NewSessionStore(r, nil)
	user, err := s.Get()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *simpleAuthHandler) authUser(session *sess.SessionStore) error {
	if err := session.Request().ParseForm(); err != nil {
		return err
	}

	userId := session.Request().FormValue("u")
	if userId == "" {
		return fmt.Errorf("userId not specified")
	}

	user, err := s.userAPI.GetUser(userId)
	if err != nil {
		return fmt.Errorf("can't get user by id %v : %v", userId, err)
	}
	return session.Save(*user)
}

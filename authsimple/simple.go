package authsimple

import (
	"context"
	"fmt"
	"net/http"
	"psychic-rat/log"
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

func (s *simpleAuthHandler) SignIn(sess *sess.SessionStore, w http.ResponseWriter) {
	if err := s.authUser(sess, w); err != nil {
		log.Logf(context.Background(), "%v", err)
		http.Error(w, "authentication failed", http.StatusForbidden)
		return
	}
}

func (s *simpleAuthHandler) Handler(w http.ResponseWriter, r *http.Request) {
	session := sess.NewSessionStore(r)
	s.SignIn(session, w)
	_, err := s.GetLoggedInUser(r)
	if err != nil {
		log.Logf(context.Background(), "%v", err)
		http.Error(w, "authentication failed", http.StatusForbidden)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (_ *simpleAuthHandler) GetLoggedInUser(r *http.Request) (*types.User, error) {
	s := sess.NewSessionStore(r)
	user, err := s.Get()
	if err != nil {
		log.Logf(context.Background(), "GetLoggedInUser: error getting session: %v", err)
	}
	return user, nil
}

func (s *simpleAuthHandler) GetUserCSRF(w http.ResponseWriter, r *http.Request) (string, error) {
	return sess.NewSessionStore(r).SetCSRF(w)
}

func (s *simpleAuthHandler) authUser(session *sess.SessionStore, w http.ResponseWriter) error {
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
	return session.Save(user, w)
}

func (s *simpleAuthHandler) LogOut(w http.ResponseWriter, r *http.Request) error {
	err := sess.NewSessionStore(r).Save(nil, w)
	if err != nil {
		return err
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

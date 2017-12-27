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
	SimpleSignIn struct {
		userAPI  UserAPI
		renderer Renderer
	}
)

func NewAuthSimple(u UserAPI, r Renderer) *SimpleSignIn {
	return &SimpleSignIn{userAPI: u, renderer: r}
}

func (s *SimpleSignIn) SignIn(sess *sess.SessionStore) {
	user, err := sess.Get()
	if err != nil {
		http.Error(sess.Writer(), err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		log.Print("no user, attempting auth")
		if err := s.authUser(sess); err != nil {
			log.Print(err)
			http.Error(sess.Writer(), "authentication failed", http.StatusForbidden)
			return
		}
	}

}

func (s *SimpleSignIn) Handler(w http.ResponseWriter, r *http.Request) {
	session := sess.NewSessionStore(r, w)
	s.SignIn(session)
	var loggedInUser types.User
	user := s.GetLoggedInUser(r)
	if user != nil {
		loggedInUser = *user
	}

	vars := struct{ User types.User }{loggedInUser}
	s.renderer.Render(w, "signin.html.tmpl", vars)
}

func (_ *SimpleSignIn) GetLoggedInUser(r *http.Request) *types.User {
	// TODO: nil is a smell. StoreReader/Writer interfaces.
	s := sess.NewSessionStore(r, nil)
	user, err := s.Get()
	if err != nil {
		log.Print(err)
		return nil
	}
	return user
}

func (s *SimpleSignIn) authUser(session *sess.SessionStore) error {
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

package auth0

import (
	"log"
	"net/http"
	"psychic-rat/sess"
	"psychic-rat/types"
)

type (
	auth0Handler struct {
		renderer    Renderer
		domain      string
		callbackURL string
		clientID    string
	}

	Renderer interface {
		Render(http.ResponseWriter, string, interface{}) error
	}
)

func NewAuth0Handler(rendr Renderer, dom string, callback string, client string) *auth0Handler {
	return &auth0Handler{
		renderer:    rendr,
		domain:      dom,
		callbackURL: callback,
		clientID:    client,
	}
}

func (a *auth0Handler) Handler(writer http.ResponseWriter, request *http.Request) {
	vars := struct {
		Auth0Domain      string
		Auth0CallbackURL string
		Auth0ClientId    string
		User             types.User
	}{
		Auth0Domain:      a.domain,
		Auth0CallbackURL: a.callbackURL,
		Auth0ClientId:    a.clientID,
	}
	a.renderer.Render(writer, "signin-auth0.html.tmpl", vars)
}

func (a *auth0Handler) GetLoggedInUser(r *http.Request) (*types.User, error) {
	s := sess.NewSessionStore(r)
	user, err := s.Get()
	if err != nil {
		log.Printf("GetLoggedInUser: error getting session: %v", err)
	}
	return user, nil
}

func (a *auth0Handler) LogOut(w http.ResponseWriter, r *http.Request) error {
	err := sess.NewSessionStore(r).Save(nil, w)
	if err != nil {
		return err
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func (a *auth0Handler) GetUserCSRF(w http.ResponseWriter, r *http.Request) (string, error) {
	return sess.NewSessionStore(r).SetCSRF(w)
}

func (a *auth0Handler) VerifyUserCSRF(r *http.Request, token string) error {
	return sess.NewSessionStore(r).VerifyCSRF(token)
}

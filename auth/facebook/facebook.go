package facebook

import (
	"fmt"
	"log"
	"net/http"
	"psychic-rat/sess"
	"psychic-rat/types"
)

const (
	authURL         string = "https://www.facebook.com/dialog/oauth"
	tokenURL        string = "https://graph.facebook.com/oauth/access_token"
	endpointProfile string = "https://graph.facebook.com/me?fields=email,first_name,last_name,link,about,id,name,picture,location"
)

type (
	handler struct {
		renderer    Renderer
		domain      string
		callbackURL string
		clientID    string
	}

	Renderer interface {
		Render(http.ResponseWriter, string, interface{}) error
	}
)

func NewHandler(rendr Renderer, callback string, clientID string) *handler {
	return &handler{
		renderer:    rendr,
		callbackURL: callback,
		clientID:    clientID,
	}
}

var state = "12345"

func (a *handler) Handler(writer http.ResponseWriter, request *http.Request) {
	vars := struct {
		AuthURL string
		User    types.User
	}{
		AuthURL: fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&state=%s", authURL, a.clientID, a.callbackURL, state),
	}
	a.renderer.Render(writer, "signin-facebook.html.tmpl", vars)
}

func (a *handler) GetLoggedInUser(r *http.Request) (*types.User, error) {
	s := sess.NewSessionStore(r)
	user, err := s.Get()
	if err != nil {
		log.Printf("GetLoggedInUser: error getting session: %v", err)
	}
	return user, nil
}

func (a *handler) LogOut(w http.ResponseWriter, r *http.Request) error {
	err := sess.NewSessionStore(r).Save(nil, w)
	if err != nil {
		return err
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func (a *handler) GetUserCSRF(w http.ResponseWriter, r *http.Request) (string, error) {
	return sess.NewSessionStore(r).SetCSRF(w)
}

func (a *handler) VerifyUserCSRF(r *http.Request, token string) error {
	return sess.NewSessionStore(r).VerifyCSRF(token)
}

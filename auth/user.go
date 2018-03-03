package auth

import (
	"fmt"
	"net/http"
	"psychic-rat/sess"
	"psychic-rat/types"
)

type userHandler struct{}

func NewUserHandler() *userHandler { return &userHandler{} }

func (u *userHandler) GetLoggedInUser(r *http.Request) (*types.User, error) {
	s := sess.NewSessionStore(r)
	user, err := s.Get()
	if err != nil {
		return nil, fmt.Errorf("GetLoggedInUser: error getting session: %v", err)
	}
	return user, nil
}

func (u *userHandler) LogOut(w http.ResponseWriter, r *http.Request) error {
	err := sess.NewSessionStore(r).Save(nil, w)
	if err != nil {
		return err
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func (u *userHandler) GetUserCSRF(w http.ResponseWriter, r *http.Request) (string, error) {
	return sess.NewSessionStore(r).SetCSRF(w)
}

func (u *userHandler) VerifyUserCSRF(r *http.Request, token string) error {
	return sess.NewSessionStore(r).VerifyCSRF(token)
}

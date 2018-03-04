package basic

import (
	"net/http"
	"psychic-rat/types"
	"strings"
)

type userHandler struct{}

func NewUserHandler() *userHandler { return &userHandler{} }

func (u *userHandler) GetLoggedInUser(r *http.Request) (*types.User, error) {
	panic("not implemented")
	username, _, ok := r.BasicAuth()
	if !ok {
		// ignore invalid headers
		return nil, nil
	}

	return &types.User{
		ID:       username,
		Email:    username + "@usermail.com",
		Fullname: username + " full",
		IsAdmin:  strings.HasPrefix(username, "admin"),
	}, nil
}

func (u *userHandler) LogOut(w http.ResponseWriter, r *http.Request) error {
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func (u *userHandler) GetUserCSRF(w http.ResponseWriter, r *http.Request) (string, error) {
	return "fakeCSRF", nil
}

func (u *userHandler) VerifyUserCSRF(r *http.Request, token string) error {
	return nil
}

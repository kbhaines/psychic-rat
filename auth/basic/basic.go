package basic

import (
	"fmt"
	"net/http"
	"psychic-rat/types"
	"strings"
)

type (
	Basic struct{}
)

func New(_ string) *Basic { return &Basic{} }

func (b *Basic) BeginAuth(w http.ResponseWriter, r *http.Request) (string, error) {
	w.Header().Set("WWW-Authenticate", `Basic realm="MY REALM"`)
	w.WriteHeader(401)
	w.Write([]byte("401 Unauthorized\n"))
	url := "/callback?p=basic"
	return url, nil
}

func (_ *Basic) Callback(w http.ResponseWriter, r *http.Request) (*types.User, error) {
	username, _, ok := r.BasicAuth()
	if !ok {
		return nil, fmt.Errorf("no authorisation header")
	}

	return &types.User{
		ID:       username,
		Email:    username + "@usermail.com",
		Fullname: username + " full",
		IsAdmin:  strings.HasPrefix(username, "admin"),
	}, nil
}

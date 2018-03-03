package basic

import (
	"encoding/base64"
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
	url := "/"
	return url, nil
}

func (_ *Basic) Callback(w http.ResponseWriter, r *http.Request) (*types.User, error) {
	header := r.Header.Get("Authorization")
	s := strings.SplitN(header, " ", 2)
	if len(s) != 2 {
		return nil, fmt.Errorf("auth header not set correctly - %s", header)
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return nil, fmt.Errorf("my auth header format wrong: %s", header)
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return nil, fmt.Errorf("yo auth header format wrong: %s", header)
	}

	return &types.User{
		ID:       pair[0],
		Email:    pair[0] + "@usermail.com",
		Fullname: pair[0],
	}, nil
}

func checkAuth(w http.ResponseWriter, r *http.Request) bool {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		return false
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return false
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return false
	}

	return pair[0] == "user" && pair[1] == "pass"
}

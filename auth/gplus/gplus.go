package gplus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"psychic-rat/types"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

const (
	authURL         string = "https://accounts.google.com/o/oauth2/auth?access_type=offline"
	tokenURL        string = "https://accounts.google.com/o/oauth2/token"
	endpointProfile string = "https://www.googleapis.com/oauth2/v2/userinfo"
)

var (
	state        = "12345"
	clientID     = os.Getenv("GPLUS_CLIENT_ID")
	clientSecret = os.Getenv("GPLUS_CLIENT_SECRET")
	cookieStore  = sessions.NewCookieStore([]byte(os.Getenv("COOKIE_KEYS")))
)

type Gplus struct {
	callbackURL string
}

func New(callbackURL string) *Gplus { return &Gplus{callbackURL} }

func (g *Gplus) BeginAuth(w http.ResponseWriter, r *http.Request) (string, error) {
	url := getOauthConf(g.callbackURL).AuthCodeURL(state)
	return url, nil
}

func getOauthConf(cb string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  cb,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}
}

func (g *Gplus) Callback(w http.ResponseWriter, r *http.Request) (*types.User, error) {
	conf := getOauthConf(g.callbackURL)

	code := r.URL.Query().Get("code")
	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("token exchange failed: %+v\n", err)
	}

	client := conf.Client(oauth2.NoContext, token)
	resp, err := client.Get(endpointProfile)
	if err != nil {
		return nil, fmt.Errorf("get profile failed: %v", err)
	}

	userProfile, err := userFromReader(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to get user from profile: %v", userProfile, err)
	}

	return userProfile, nil
}

func userFromReader(reader io.Reader) (*types.User, error) {
	u := struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		About     string `json:"about"`
		Name      string `json:"name"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Link      string `json:"link"`
		Location  struct {
			Name string `json:"name"`
		} `json:"location"`
	}{}

	err := json.NewDecoder(reader).Decode(&u)
	if err != nil {
		return nil, err
	}

	user := &types.User{}
	user.FirstName = u.FirstName
	user.Fullname = u.Name
	user.Email = u.Email
	user.ID = "gplus" + u.ID
	user.Country = u.Location.Name

	return user, err
}

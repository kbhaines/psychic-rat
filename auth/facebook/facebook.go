package facebook

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
	authURL         string = "https://www.facebook.com/dialog/oauth"
	tokenURL        string = "https://graph.facebook.com/oauth/access_token"
	endpointProfile string = "https://graph.facebook.com/me?fields=email,first_name,last_name,id,name"
	facebookCookie  string = "fb"
)

var (
	clientID     = os.Getenv("FACEBOOK_CLIENT_ID")
	clientSecret = os.Getenv("FACEBOOK_CLIENT_SECRET")
	cookieStore  = sessions.NewCookieStore([]byte(os.Getenv("COOKIE_KEYS")))
)

type (
	Facebook struct {
		callbackURL string
		csrf        CSRFValidator
	}

	CSRFValidator interface {
		GetUserCSRF(w http.ResponseWriter, r *http.Request) (string, error)
		VerifyUserCSRF(r *http.Request, token string) error
	}
)

func New(callbackURL string, csrfValidator CSRFValidator) *Facebook {
	return &Facebook{callbackURL, csrfValidator}
}

func (f *Facebook) BeginAuth(w http.ResponseWriter, r *http.Request) (string, error) {
	state, err := f.csrf.GetUserCSRF(w, r)
	if err != nil {
		return "", fmt.Errorf("BeginAuth: CSRF error: %v", err)
	}
	url := getOauthConf(f.callbackURL).AuthCodeURL(state)
	return url, nil
}

func getOauthConf(cb string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  cb,
		Scopes:       []string{"public_profile", "email", "user_hometown"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}
}

func (f *Facebook) Callback(w http.ResponseWriter, r *http.Request) (*types.User, error) {
	conf := getOauthConf(f.callbackURL)

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	err := f.csrf.VerifyUserCSRF(r, state)
	if err != nil {
		return nil, fmt.Errorf("Callback: state validation error: %v", err)
	}

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
		Picture   struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
		Location struct {
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
	user.ID = "facebook" + u.ID
	user.Country = u.Location.Name

	return user, err
}

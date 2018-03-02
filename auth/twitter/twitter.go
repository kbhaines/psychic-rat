package twitter

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"psychic-rat/log"
	"psychic-rat/types"

	"github.com/gorilla/sessions"
	"github.com/mrjones/oauth"
)

const (
	requestURL      = "https://api.twitter.com/oauth/request_token"
	authorizeURL    = "https://api.twitter.com/oauth/authorize"
	authenticateURL = "https://api.twitter.com/oauth/authenticate"
	tokenURL        = "https://api.twitter.com/oauth/access_token"
	endpointProfile = "https://api.twitter.com/1.1/account/verify_credentials.json"
	tokenCookie     = "oauth-twitter"
)

type (
	handler struct {
		callbackURL  string
		clientID     string
		clientSecret string
	}

	Twitter struct {
		callbackURL string
	}
)

var (
	clientID     = os.Getenv("TWITTER_CLIENT_ID")
	clientSecret = os.Getenv("TWITTER_CLIENT_SECRET")
	cookieStore  = sessions.NewCookieStore([]byte(os.Getenv("COOKIE_KEYS")))
)

func init() {
	gob.Register(&oauth.RequestToken{})
}

func New(callbackURL string) *Twitter { return &Twitter{callbackURL} }

func (t *Twitter) BeginAuth(w http.ResponseWriter, r *http.Request) (string, error) {
	consumer := newConsumer()
	requestToken, url, err := consumer.GetRequestTokenAndUrl(t.callbackURL)
	if err != nil {
		return "", fmt.Errorf("unable to get token from twitter: %v", err)
	}

	session, err := cookieStore.Get(r, tokenCookie)
	session.Values["token"] = *requestToken
	err = session.Save(r, w)
	if err != nil {
		return "", fmt.Errorf("unable to save session: %v", err)
	}
	log.Logf(r.Context(), "stored token: %v", requestToken)
	return url, nil
}

func newConsumer() *oauth.Consumer {
	return oauth.NewConsumer(
		clientID,
		clientSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   requestURL,
			AuthorizeTokenUrl: authorizeURL,
			AccessTokenUrl:    tokenURL,
		})
}

func (_ *Twitter) Callback(w http.ResponseWriter, r *http.Request) (*types.User, error) {
	sessv, err := cookieStore.Get(r, tokenCookie)
	if err != nil {
		return nil, fmt.Errorf("unable to get twitter cookie from user session: %v", err)
	}
	token, tokenOK := sessv.Values["token"].(*oauth.RequestToken)
	if !tokenOK {
		return nil, fmt.Errorf("unable to get token (%s)", token)
	}

	values := r.URL.Query()
	verificationCode := values.Get("oauth_verifier")
	//tokenKey := values.Get("oauth_token")

	//log.Logf(r.Context(), "token: %s, tokenKey: %s", token.Token, tokenKey)

	c := newConsumer()
	accessToken, err := c.AuthorizeToken(token, verificationCode)
	if err != nil {
		return nil, fmt.Errorf("unable to authorize token: %v", err)
	}

	/*client, err := c.MakeHttpClient(accessToken)
	if err != nil {
		log.Errorf(r.Context(), "unable to make client: %v", err)
		http.Error(w, "token process error", http.StatusInternalServerError)
		return
	}
	*/

	// todo: Get is deprecated
	response, err := c.Get(
		endpointProfile,
		map[string]string{"include_entities": "false", "skip_status": "true", "include_email": "true"},
		accessToken)
	if err != nil {
		return nil, fmt.Errorf("profile error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("profile error: %v", err)
	}
	userProfile, err := userFromReader(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to create a user %v :%v", userProfile, err)
	}
	return userProfile, nil
}

func userFromReader(reader io.Reader) (*types.User, error) {
	u := struct {
		ID       string `json:"id_str"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		Location string `json:"location"`
	}{}

	err := json.NewDecoder(reader).Decode(&u)
	if err != nil {
		return nil, err
	}

	user := &types.User{}
	user.Fullname = u.Name
	user.Email = u.Email
	user.ID = "twitter" + u.ID
	user.Country = u.Location

	return user, err
}

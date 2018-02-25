package twitter

import (
	"encoding/gob"
	"encoding/json"
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

type handler struct {
	callbackURL  string
	clientID     string
	clientSecret string
}

var (
	clientID     = os.Getenv("TWITTER_CLIENT_ID")
	clientSecret = os.Getenv("TWITTER_CLIENT_SECRET")
	cookieStore  = sessions.NewCookieStore([]byte(os.Getenv("COOKIE_KEYS")))
)

func init() {
	gob.Register(&oauth.RequestToken{})
}

func BeginAuth(w http.ResponseWriter, r *http.Request) {
	consumer := newConsumer()
	requestToken, url, err := consumer.GetRequestTokenAndUrl("http://localhost:8080/auth/twitter/callback")
	if err != nil {
		log.Errorf(r.Context(), "unable to get token from twitter: %v", err)
		http.Error(w, "twitter token error", http.StatusInternalServerError)
		return
	}

	session, err := cookieStore.Get(r, tokenCookie)
	session.Values["token"] = *requestToken
	err = session.Save(r, w)
	if err != nil {
		log.Errorf(r.Context(), "unable to save session: %v", err)
		http.Error(w, "twitter token error", http.StatusInternalServerError)
		return
	}

	log.Logf(r.Context(), "stored token: %v", requestToken)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	return
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := cookieStore.Get(r, tokenCookie)
	if err != nil {
		log.Errorf(r.Context(), "unable to get twitter cookie from user session: %v", err)
		http.Error(w, "twitter token error", http.StatusInternalServerError)
		return
	}
	log.Logf(r.Context(), "session: %+v", sess.Values)
	token, tokenOK := sess.Values["token"].(*oauth.RequestToken)
	if !tokenOK {
		log.Errorf(r.Context(), "unable to get token (%s)", token)
		http.Error(w, "token process error", http.StatusInternalServerError)
		return
	}

	values := r.URL.Query()
	verificationCode := values.Get("oauth_verifier")
	tokenKey := values.Get("oauth_token")

	log.Logf(r.Context(), "token: %s, tokenKey: %s", token.Token, tokenKey)

	c := newConsumer()
	accessToken, err := c.AuthorizeToken(token, verificationCode)
	if err != nil {
		log.Errorf(r.Context(), "unable to authorize token: %v", err)
		http.Error(w, "token process error", http.StatusInternalServerError)
		return
	}

	/*client, err := c.MakeHttpClient(accessToken)
	if err != nil {
		log.Errorf(r.Context(), "unable to make client: %v", err)
		http.Error(w, "token process error", http.StatusInternalServerError)
		return
	}
	*/

	response, err := c.Get(
		endpointProfile,
		map[string]string{"include_entities": "false", "skip_status": "true", "include_email": "true"},
		accessToken)
	if err != nil {
		log.Errorf(r.Context(), "profile error: %v", err)
		http.Error(w, "token process error", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Errorf(r.Context(), "profile error: %v", err)
		http.Error(w, "token process error", http.StatusInternalServerError)
		return
	}
	user, err := userFromReader(response.Body)
	log.Logf(r.Context(), "user: %v", user)
	return
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

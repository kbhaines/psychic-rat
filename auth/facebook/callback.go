package facebook

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"psychic-rat/log"
	"psychic-rat/sess"
	"psychic-rat/types"

	"golang.org/x/oauth2"
)

type UserAPI interface {
	GetUser(id string) (*types.User, error)
	AddUser(types.User) error
}

var (
	userAPI UserAPI
)

func Init(a UserAPI) {
	userAPI = a
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
		ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/facebook",
		//Scopes:       []string{"openid", "profile", "user_metadata"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}

	code := r.URL.Query().Get("code")

	token, err := conf.Exchange(oauth2.NoContext, code)
	log.Logf(r.Context(), "token = %+v\n", token)
	if err != nil {
		log.Logf(r.Context(), "token exchange failed: %+v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := conf.Client(oauth2.NoContext, token)
	resp, err := client.Get(endpointProfile)
	if err != nil {
		log.Errorf(r.Context(), "get profile failed: %+v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userProfile, err := userFromReader(resp.Body)
	defer resp.Body.Close()

	userRecord, err := addUserIfNotExists(userProfile)
	if err != nil {
		log.Errorf(r.Context(), "unable to create a user %v :%v", userRecord, err)
		return
	}
	err = sess.NewSessionStore(r).Save(userRecord, w)
	if err != nil {
		log.Errorf(r.Context(), "unable to save user into session: %v", err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func addUserIfNotExists(u *types.User) (*types.User, error) {
	existing, err := userAPI.GetUser(u.ID)
	if err != nil {
		return u, userAPI.AddUser(*u)
	}
	return existing, nil
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
package auth

import (
	"net/http"
	"psychic-rat/log"
	"psychic-rat/sess"
	"psychic-rat/types"
	"sort"
)

type (
	UserAPI interface {
		GetUser(id string) (*types.User, error)
		AddUser(types.User) error
	}

	AuthHandler interface {
		BeginAuth(w http.ResponseWriter, r *http.Request) (string, error)
		Callback(w http.ResponseWriter, r *http.Request) (*types.User, error)
	}
)

var (
	authProviders map[string]AuthHandler
	userAPI       UserAPI
)

func Init(u UserAPI, ap map[string]AuthHandler) {
	userAPI = u
	authProviders = ap
}

func AuthInit(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("p")
	handler, ok := authProviders[provider]
	if !ok {
		log.Errorf(r.Context(), "provider %s not found", provider)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	redirect, err := handler.BeginAuth(w, r)
	if err != nil {
		log.Errorf(r.Context(), "unable to start auth process for %s: %v", provider, err)
		http.Error(w, "could not start auth", http.StatusInternalServerError)
		return
	}
	country := r.URL.Query().Get("c")
	if !validCountry(country) {
		log.Errorf(r.Context(), "invalid country: %s", country)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	user := types.User{Country: country}
	sess.NewSessionStore(r).Save(&user, w)
	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("p")
	handler, ok := authProviders[provider]
	if !ok {
		log.Errorf(r.Context(), "callback provider %s not found", provider)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	country := r.URL.Query().Get("c")
	if country == "" {
		user, err := sess.NewSessionStore(r).Get()
		if err != nil || user == nil {
			log.Errorf(r.Context(), "could not read user back to get country: %v", err)
			http.Error(w, "auth error has occurred", http.StatusInternalServerError)
			return
		}
		country = user.Country
	}
	if !validCountry(country) {
		log.Errorf(r.Context(), "invalid country in callback: %s", country)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := handler.Callback(w, r)
	if err != nil {
		log.Errorf(r.Context(), "could not handle callback: %v", err)
		http.Error(w, "auth error has occurred", http.StatusInternalServerError)
		return
	}
	user.Country = country

	user, err = addUserIfNotExists(user)
	if err != nil {
		log.Errorf(r.Context(), "unable to create a user %v :%v", user, err)
		http.Error(w, "auth error has occurred", http.StatusInternalServerError)
		return
	}

	err = sess.NewSessionStore(r).Save(user, w)
	if err != nil {
		log.Errorf(r.Context(), "unable to save user into session: %v", err)
		http.Error(w, "auth error has occurred", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func addUserIfNotExists(u *types.User) (*types.User, error) {
	existing, err := userAPI.GetUser(u.ID)
	if err != nil {
		return u, userAPI.AddUser(*u)
	}
	return existing, nil
}

func validCountry(c string) bool {
	idx := sort.SearchStrings(countries, c)
	return idx < len(countries) && countries[idx] == c
}

var countries = []string{
	"AD", "AE", "AF", "AG", "AI", "AL", "AM", "AN", "AO", "AQ", "AR", "AS",
	"AT", "AU", "AW", "AX", "AZ", "BA", "BB", "BD", "BE", "BF", "BG", "BH",
	"BI", "BJ", "BM", "BN", "BO", "BR", "BS", "BT", "BV", "BW", "BY", "BZ",
	"CA", "CC", "CD", "CF", "CG", "CH", "CI", "CK", "CL", "CM", "CN", "CO",
	"CR", "CU", "CV", "CX", "CY", "CZ", "DE", "DJ", "DK", "DM", "DO", "DZ",
	"EC", "EE", "EG", "EH", "ER", "ES", "ET", "FI", "FJ", "FK", "FM", "FO",
	"FR", "GA", "GB", "GD", "GE", "GF", "GG", "GH", "GI", "GL", "GM", "GN",
	"GP", "GQ", "GR", "GS", "GT", "GU", "GW", "GY", "HK", "HM", "HN", "HR",
	"HT", "HU", "ID", "IE", "IL", "IM", "IN", "IO", "IQ", "IR", "IS", "IT",
	"JE", "JM", "JO", "JP", "KE", "KG", "KH", "KI", "KM", "KN", "KP", "KR",
	"KW", "KY", "KZ", "LA", "LB", "LC", "LI", "LK", "LR", "LS", "LT", "LU",
	"LV", "LY", "MA", "MC", "MD", "ME", "MG", "MH", "MK", "ML", "MM", "MN",
	"MO", "MP", "MQ", "MR", "MS", "MT", "MU", "MV", "MW", "MX", "MY", "MZ",
	"NA", "NC", "NE", "NF", "NG", "NI", "NL", "NO", "NP", "NR", "NU", "NZ",
	"OM", "PA", "PE", "PF", "PG", "PH", "PK", "PL", "PM", "PN", "PR", "PS",
	"PT", "PW", "PY", "QA", "RE", "RO", "RS", "RU", "RW", "SA", "SB", "SC",
	"SD", "SE", "SG", "SH", "SI", "SJ", "SK", "SL", "SM", "SN", "SO", "SR",
	"ST", "SV", "SY", "SZ", "TC", "TD", "TF", "TG", "TH", "TJ", "TK", "TL",
	"TM", "TN", "TO", "TR", "TT", "TV", "TW", "TZ", "UA", "UG", "UM", "US",
	"UY", "UZ", "VA", "VC", "VE", "VG", "VI", "VN", "VU", "WF", "WS", "YE",
	"YT", "ZA", "ZM", "ZW",
}

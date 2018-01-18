package web

import (
	"net/http"
	"psychic-rat/auth0"
	"psychic-rat/web/admin"
	"psychic-rat/web/dispatch"
	"psychic-rat/web/pub"
)

const (
	HomePage      = "/"
	SignInPage    = "/signin"
	SignOutPage   = "/signout"
	PledgePage    = "/pledge"
	ThanksPage    = "/thanks"
	NewItem       = "/newitem"
	AdminNewItems = "/admin/newitems"
	Callback      = "/callback"
)

var (
	uriHandlers = []dispatch.URIHandler{
		{HomePage, pub.HomePageHandler},
		{SignInPage, pub.SignInPageHandler},
		{SignOutPage, pub.SignOutPageHandler},
		{PledgePage, pub.PledgePageHandler},
		{NewItem, pub.NewItemHandler},
		{ThanksPage, pub.ThanksPageHandler},
		{Callback, auth0.CallbackHandler},
		{AdminNewItems, admin.AdminItemHandler},
	}

	flags struct {
		enableAuth0, sqldb bool
	}
)

func handlerForDirs(mux *http.ServeMux, dir ...string) {
	for _, d := range dir {
		mux.Handle("/"+d+"/", http.StripPrefix("/"+d, http.FileServer(http.Dir(d+"/"))))
	}
}

func Handler() http.Handler {
	hmux := http.NewServeMux()
	handlerForDirs(hmux, "css", "js", "images")
	for _, h := range uriHandlers {
		hmux.HandleFunc(h.URI, h.Handler)
	}
	return hmux
}

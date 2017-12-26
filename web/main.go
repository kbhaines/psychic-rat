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

func Handler() http.Handler {
	hmux := http.NewServeMux()
	for _, h := range uriHandlers {
		hmux.HandleFunc(h.URI, h.Handler)
	}
	return hmux
}

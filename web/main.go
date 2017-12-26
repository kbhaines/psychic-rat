package web

import (
	"net/http"
	"psychic-rat/auth0"
	"psychic-rat/types"
	"psychic-rat/web/admin"
	"psychic-rat/web/dispatch"
)

type (
	APIS struct {
		Item    ItemAPI
		NewItem NewItemAPI
		Pledge  PledgeAPI
		User    UserAPI
	}

	// TODO: lots of interfaces here and they need to be split into smaller
	// ones, along with splitting the web module as well.

	// TODO: APIs need consistent parameter and return style

	ItemAPI interface {
		ListItems() ([]types.Item, error)
		GetItem(id int) (types.Item, error)
	}

	NewItemAPI interface {
		AddNewItem(types.NewItem) (*types.NewItem, error)
	}

	PledgeAPI interface {
		AddPledge(itemId int, userId string) (int, error)
	}

	UserAPI interface {
		GetUser(userId string) (*types.User, error)
	}

	NewItemPost struct {
		Company string
		Make    string
		Model   string
	}
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
		{HomePage, HomePageHandler},
		{SignInPage, SignInPageHandler},
		{PledgePage, PledgePageHandler},
		{NewItem, NewItemHandler},
		{ThanksPage, ThanksPageHandler},
		{Callback, auth0.CallbackHandler},
		{AdminNewItems, admin.AdminItemHandler},
	}

	apis APIS

	flags struct {
		enableAuth0, sqldb bool
	}
)

func Init(a APIS) {
	apis = a
}

func Handler() http.Handler {
	hmux := http.NewServeMux()
	for _, h := range uriHandlers {
		hmux.HandleFunc(h.URI, h.Handler)
	}
	return hmux
}

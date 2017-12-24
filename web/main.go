package web

import (
	"net/http"
	"psychic-rat/auth0"
	"psychic-rat/types"
	"psychic-rat/web/admin"
	"psychic-rat/web/dispatch"
)

type (

	// TODO: Need for this should go when package is split/refactored
	API struct {
		Company CompanyAPI
		Item    ItemAPI
		NewItem NewItemAPI
		Pledge  PledgeAPI
		User    UserAPI
	}

	// TODO: lots of interfaces here and they need to be split into smaller
	// ones, along with splitting the web module as well.

	// TODO: APIs need consistent parameter and return style

	CompanyAPI interface {
		GetCompanies() ([]types.Company, error)
		GetCompany(int) (types.Company, error)
		AddCompany(co types.Company) (*types.Company, error)
	}

	ItemAPI interface {
		ListItems() ([]types.Item, error)
		AddItem(types.Item) (*types.Item, error)
		GetItem(id int) (types.Item, error)
	}

	NewItemAPI interface {
		AddNewItem(types.NewItem) (*types.NewItem, error)
		ListNewItems() ([]types.NewItem, error)
		DeleteNewItem(int) error
	}

	AdminApi interface {
		ApproveItem(id int) error
	}

	PledgeAPI interface {
		AddPledge(itemId int, userId string) (int, error)
		//ListPledges() PledgeListing
	}

	UserAPI interface {
		GetUser(userId string) (*types.User, error)
		AddUser(types.User) error
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

	apis API

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

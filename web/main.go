package web

import (
	"net/http"
	"psychic-rat/api/rest"
	"psychic-rat/auth0"
	"psychic-rat/mdl"
	"psychic-rat/sqldb"
	"psychic-rat/types"
)

type (
	URIHandler struct {
		URI     string
		Handler http.HandlerFunc
	}

	API struct {
		Company CompanyAPI
		Item    ItemAPI
		Pledge  PledgeAPI
		User    UserAPI
	}

	CompanyAPI interface {
		GetCompanies() (types.CompanyListing, error)
		GetById(mdl.ID) (types.CompanyElement, error)
	}

	ItemAPI interface {
		ListItems() (types.ItemReport, error)
		GetById(id mdl.ID) (types.ItemElement, error)
		AddItem(item mdl.NewItem) error
		ListNewItems() ([]mdl.NewItem, error)
		ApproveItem(id mdl.ID) error
	}

	PledgeAPI interface {
		NewPledge(itemId mdl.ID, userId mdl.ID) (mdl.ID, error)
		//ListPledges() PledgeListing
	}

	UserAPI interface {
		GetById(userId mdl.ID) (*mdl.User, error)
		Create(mdl.User) error
	}
)

var (
	uriHandlers = []URIHandler{
		{rest.HomePage, HomePageHandler},
		{rest.SignInPage, SignInPageHandler},
		{rest.PledgePage, PledgePageHandler},
		{rest.ThanksPage, ThanksPageHandler},
		{"/callback", auth0.CallbackHandler},
	}

	apis API
	db   sqldb.DB

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

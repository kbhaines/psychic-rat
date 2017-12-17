package web

import (
	"net/http"
	"psychic-rat/api/rest"
	"psychic-rat/auth0"
	"psychic-rat/mdl"
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
		NewItem NewItemAPI
		Pledge  PledgeAPI
		User    UserAPI
	}

	CompanyAPI interface {
		GetCompanies() ([]types.Company, error)
		GetCompany(int) (types.Company, error)
		NewCompany(co types.Company) (*types.Company, error)
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
		NewPledge(itemId int, userId string) (int, error)
		//ListPledges() PledgeListing
	}

	UserAPI interface {
		GetUser(userId string) (*mdl.User, error)
		CreateUser(mdl.User) error
	}

	NewItemPost struct {
		Company string
		Make    string
		Model   string
	}

	NewItemAdminPost struct {
		ID          int
		Add         bool
		Delete      bool
		Pledge      bool
		ItemID      int
		CompanyID   int
		UserID      string
		UserCompany string
		UserMake    string
		UserModel   string
	}
)

var (
	uriHandlers = []URIHandler{
		{rest.HomePage, HomePageHandler},
		{rest.SignInPage, SignInPageHandler},
		{rest.PledgePage, PledgePageHandler},
		{rest.NewItem, NewItemHandler},
		{rest.ThanksPage, ThanksPageHandler},
		{"/callback", auth0.CallbackHandler},

		{"/admin/newitems", AdminItemHandler},
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

package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"psychic-rat/api/rest"
	"psychic-rat/auth0"
	"psychic-rat/mdl"
	"psychic-rat/types"
)

type (
	UriHandler struct {
		URI     string
		Handler http.HandlerFunc
	}

	API struct {
		Company CompanyApi
		Item    ItemApi
		Pledge  PledgeApi
		User    UserApi
	}

	CompanyApi interface {
		GetCompanies() (types.CompanyListing, error)
		GetById(mdl.ID) (types.CompanyElement, error)
	}

	ItemApi interface {
		ListItems() (types.ItemReport, error)
		GetById(id mdl.ID) (types.ItemElement, error)
		AddItem(item mdl.NewItem) error
		ListNewItems() ([]mdl.NewItem, error)
		ApproveItem(id mdl.ID) error
	}

	PledgeApi interface {
		NewPledge(itemId mdl.ID, userId mdl.ID) (mdl.ID, error)
		//ListPledges() PledgeListing
	}

	UserApi interface {
		GetById(userId mdl.ID) (*mdl.User, error)
		Create(mdl.User) error
	}
)

var (
	uriHandlers = []UriHandler{
		{rest.HomePage, HomePageHandler},
		{rest.SignInPage, SignInPageHandler},
		{rest.PledgePage, PledgePageHandler},
		{rest.ThanksPage, ThanksPageHandler},
		{"/callback", auth0.CallbackHandler},
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

func ToJson(writer io.Writer, v interface{}) {
	fmt.Fprintf(writer, "%s", ToJsonString(v))
}

func ToJsonString(v interface{}) string {
	js, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("unable to convert %T (%v)to json", v, v))
	}
	return string(js)
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"psychic-rat/api/rest"
	"psychic-rat/auth0"
	"psychic-rat/mdl"
	"psychic-rat/types"

	"github.com/gorilla/context"
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
	UriHandlers = []UriHandler{
		{rest.CompanyApi, CompanyHandler},
		{rest.ItemApi, ItemHandler},
		{rest.PledgeApi, PledgeHandler},

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

func handler() http.Handler {
	hmux := http.NewServeMux()
	for _, h := range UriHandlers {
		hmux.HandleFunc(h.URI, h.Handler)
	}
	return hmux
}

func main() {
	flag.BoolVar(&flags.enableAuth0, "auth0", false, "enable auth0 function")
	flag.BoolVar(&flags.sqldb, "sqldb", false, "enable real database")
	flag.Parse()

	http.ListenAndServe("localhost:8080", context.ClearHandler(handler()))
}

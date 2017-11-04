package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"psychic-rat/api/rest"
	"psychic-rat/mdl"
	"psychic-rat/types"

	"github.com/gorilla/context"
)

type (
	UriHandler struct {
		Uri     string
		Handler http.HandlerFunc
	}

	Api struct {
		Company CompanyApi
		Item    ItemApi
		Pledge  PledgeApi
		User    UserApi
	}

	CompanyApi interface {
		GetCompanies() (types.CompanyListing, error)
		GetById(mdl.Id) (types.CompanyElement, error)
	}

	ItemApi interface {
		//AddItem(make string, model string, company company.Id) error
		ListItems() (types.ItemReport, error)
		GetById(id mdl.Id) (types.ItemElement, error)
	}

	PledgeApi interface {
		NewPledge(itemId mdl.Id, userId mdl.Id) (mdl.Id, error)
		//ListPledges() PledgeListing
	}

	UserApi interface {
		GetById(userId mdl.Id) (*mdl.UserRecord, error)
		Create(mdl.UserRecord) error
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
		{"/callback", CallbackHandler},
	}

	apis Api
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
		hmux.HandleFunc(h.Uri, h.Handler)
	}
	return hmux
}

func main() {
	http.ListenAndServe("localhost:8080", context.ClearHandler(handler()))
}

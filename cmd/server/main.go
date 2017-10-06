package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"psychic-rat/api"
	"psychic-rat/api/rest"
	"psychic-rat/impl"
	"psychic-rat/mdl"

	"github.com/gorilla/context"
)

type UriHandler struct {
	Uri     string
	Handler http.HandlerFunc
}

var (
	UriHandlers = []UriHandler{
		{rest.CompanyApi, CompanyHandler},
		{rest.ItemApi, ItemHandler},
		{rest.PledgeApi, PledgeHandler},

		{rest.HomePage, HomePageHandler},
		{rest.SignInPage, SignInPageHandler},
		{rest.PledgePage, PledgePageHandler},
		{rest.ThanksPage, ThanksPageHandler},
	}

	apis api.Api
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

func setupMockData() {
	companies := impl.GetRepos().Company
	companies.Create(mdl.NewCompany(mdl.Id("1"), "bigco1"))
	companies.Create(mdl.NewCompany(mdl.Id("2"), "bigco2"))
	companies.Create(mdl.NewCompany(mdl.Id("3"), "bigco3"))

	items := impl.GetRepos().Item
	items.Create(mdl.NewItem("phone", "abc", mdl.Id("1")))
	items.Create(mdl.NewItem("phone", "xyz", mdl.Id("1")))
	items.Create(mdl.NewItem("tablet", "gt1", mdl.Id("1")))
	items.Create(mdl.NewItem("tablet", "tab4", mdl.Id("2")))
	items.Create(mdl.NewItem("tablet", "tab8", mdl.Id("2")))

	users := impl.GetRepos().User
	users.Create(mdl.UserRecord{Id: mdl.Id("1234"), Email: "kbh@x.y", FirstName: "Kevin"})
}

func init() {
}

func main() {
	apis = impl.GetApi()
	setupMockData()
	for _, h := range UriHandlers {
		http.HandleFunc(h.Uri, h.Handler)
	}

	http.ListenAndServe("localhost:8080", context.ClearHandler(http.DefaultServeMux))
}

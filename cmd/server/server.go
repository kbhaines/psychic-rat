package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"psychic-rat/api/rest"
	"psychic-rat/mdl/company"
	"psychic-rat/mdl/item"
	"psychic-rat/repo/companyrepo"
	"psychic-rat/repo/itemrepo"
)

type UriHandler struct {
	Uri     string
	Handler http.HandlerFunc
}

var UriHandlers = []UriHandler{
	{rest.CompanyApi, CompanyHandler},
	{rest.ItemApi, ItemHandler},
	{rest.PledgeApi, PledgeHandler},

	{rest.HomePage, HomePageHandler},
	{rest.SignInPage, SignInPageHandler},
	{rest.PledgePage, PledgePageHandler},
	{rest.ThanksPage, ThanksPageHandler},
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

func main() {
	companies := companyrepo.GetCompanyRepoMapImpl()
	companies.Create(company.New(company.Id("1"), "bigco1"))
	companies.Create(company.New(company.Id("2"), "bigco2"))
	companies.Create(company.New(company.Id("3"), "bigco3"))

	items := itemrepo.GetItemRepoMapImpl()
	items.Create(item.New("phone", "abc", company.Id("1")))
	items.Create(item.New("phone", "xyz", company.Id("1")))
	items.Create(item.New("tablet", "gt1", company.Id("1")))
	items.Create(item.New("tablet", "tab4", company.Id("2")))
	items.Create(item.New("tablet", "tab8", company.Id("2")))

	for _, h := range UriHandlers {
		http.HandleFunc(h.Uri, h.Handler)
	}

	http.ListenAndServe("localhost:8080", nil)
}

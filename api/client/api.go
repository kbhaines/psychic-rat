package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	ap "psychic-rat/api"
	"psychic-rat/mdl/company"
	"psychic-rat/mdl/item"
	"strings"
)

type Api interface {
	GetCompanies() (ap.CompanyListing, error)
	GetItems(companyId company.Id) (ap.ItemReport, error)
	NewPledge(itemId item.Id) (ap.PledgeListing, error)
}

type api struct {
	url url.URL
}

func New(targetUrl string) Api {
	u, err := url.Parse(targetUrl)
	if err != nil {
		panic(fmt.Errorf("invalid targetUrl: %v - %v", targetUrl, err))
	}
	return &api{url: *u}
}

func (a *api) doGet() (http.Response, error) {
	panic("not implemented")
}

func (a *api) GetCompanies() (ap.CompanyListing, error) {
	resp, err := http.Get(a.url.String() + ap.CompanyApi)
	if err != nil {
		return ap.CompanyListing{}, fmt.Errorf("get companies failed: %v", err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ap.CompanyListing{}, fmt.Errorf("error reading conpany response: %v", err)
	}
	return ap.CompaniesFromJson(bytes)
}

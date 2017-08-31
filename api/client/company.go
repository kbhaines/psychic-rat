package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	api2 "psychic-rat/api"
	"psychic-rat/mdl/company"
	"psychic-rat/mdl/item"
)

type Api interface {
	GetCompanies() ([]company.Record, error)
	GetItems(companyId company.Id) ([]item.Record, error)
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

func (a *api) GetCompanies() ([]company.Record, error) {
	resp, err := http.Get(a.url.String() + api2.CompanyApi)
	if err != nil {
		return nil, fmt.Errorf("get companies failed: %v", err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading conpany response: %v", err)
	}
	return api2.CompaniesFromJson(bytes)
}

func (a api) GetItems(id company.Id) ([]item.Record, error) {
	resp, err := http.Get(a.url.String() + api2.ItemApi + fmt.Sprintf("?company=%v", id))
	if err != nil {
		return nil, fmt.Errorf("get items failed: %v", err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading item response: %v", err)
	}
	return api2.ItemsFromJson(bytes)
}

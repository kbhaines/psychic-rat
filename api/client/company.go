package client

import (
	"psychic-rat/mdl/company"
	"net/url"
	"fmt"
	"net/http"
	api2 "psychic-rat/api"
	"io/ioutil"
)

type Api interface {
	GetCompanies() ([]company.Record, error)
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

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

func (a *api) GetItems(id company.Id) (ap.ItemReport, error) {
	resp, err := http.Get(a.url.String() + ap.ItemApi + fmt.Sprintf("?company=%v", id))
	report := ap.ItemReport{}
	if err != nil {
		return report, fmt.Errorf("get items failed: %v", err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return report, fmt.Errorf("error reading item response: %v", err)
	}
	return ap.ItemsFromJson(bytes)
}

func (a *api) NewPledge(itemId item.Id) (ap.PledgeListing, error) {
	jsonString := ap.ToJsonString(ap.NewPledgeRequest(itemId))
	body := strings.NewReader(jsonString)
	resp, err := http.Post(a.url.String()+ap.PledgeApi, "application/json", body)
	listing := ap.PledgeListing{}
	if resp.StatusCode != http.StatusOK {
		return listing, fmt.Errorf("request returned: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bytes, &listing)
	if err != nil {
		return listing, fmt.Errorf("unable to decode response: %v", err)
	}
	return listing, nil
}

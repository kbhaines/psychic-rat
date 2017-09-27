package api

import (
	"psychic-rat/mdl"
	"time"
)

// Companies
type CompanyApi interface {
	GetCompanies() (CompanyListing, error)
	GetById(mdl.Id) (CompanyElement, error)
}

type CompanyListing struct {
	Companies []CompanyElement `json:"companies"`
}

type CompanyElement struct {
	Id   mdl.Id `json:"id"`
	Name string `json:"name"`
}

// Items
type ItemApi interface {
	//AddItem(make string, model string, company company.Id) error
	ListItems() (ItemReport, error)
	GetById(id mdl.Id) (ItemElement, error)
}

type ItemReport struct {
	Items []ItemElement `json:"items"`
}

type ItemElement struct {
	Id      mdl.Id `json:"id"`
	Make    string `json:"make"`
	Model   string `json:"model"`
	Company string `json:"company"`
}

// Pledges
type PledgeApi interface {
	NewPledge(itemId mdl.Id, userId mdl.Id) (mdl.Id, error)
	//ListPledges() PledgeListing
}

type PledgeListing struct {
	Pledges []PledgeElement `json:"pledges"`
}

type PledgeElement struct {
	PledgeId  mdl.Id      `json:"id"`
	Item      ItemElement `json:"item"`
	Timestamp time.Time   `json:"timestamp"`
}

type PledgeRequest struct {
	ItemId mdl.Id `json:"itemId"`
}

type Api struct {
	Company CompanyApi
	Item    ItemApi
	Pledge  PledgeApi
}

func init() {
	api = Api{
		Company: getRepoCompanyApi(),
		Item:    getRepoItemApi(),
		Pledge:  getRepoPledgeApi(),
	}
}

var api Api

func Get() Api {
	return api
}

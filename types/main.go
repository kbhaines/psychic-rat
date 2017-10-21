package types

import (
	"psychic-rat/mdl"
	"time"
)

type CompanyListing struct {
	Companies []CompanyElement `json:"companies"`
}

type CompanyElement struct {
	Id   mdl.Id `json:"id"`
	Name string `json:"name"`
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

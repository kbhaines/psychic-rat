package types

import (
	"time"
)

type CompanyListing struct {
	Companies []Company `json:"companies"`
}

type Company struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ItemReport struct {
	Items []Item `json:"items"`
}

type Item struct {
	Id      int    `json:"id"`
	Make    string `json:"make"`
	Model   string `json:"model"`
	Company string `json:"company"`
}

type PledgeListing struct {
	Pledges []Pledge `json:"pledges"`
}

type Pledge struct {
	PledgeId  int       `json:"id"`
	Item      Item      `json:"item"`
	Timestamp time.Time `json:"timestamp"`
}

type PledgeRequest struct {
	ItemId int `json:"itemId"`
}

// NewItem is a record of a (user) request to add a new item to the database of
// Items. 'Company' will be used if also adding in a new company to the
// database. IsPledge is true if the user also pledged the item
type NewItem struct {
	Id        int
	UserID    int
	IsPledge  bool
	Make      string
	Model     string
	Company   string
	CompanyID int
	Timestamp time.Time
}

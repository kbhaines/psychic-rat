package types

import (
	"time"
)

type Company struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Item struct {
	ID      int    `json:"id"`
	Make    string `json:"make"`
	Model   string `json:"model"`
	Company Company
}

type Pledge struct {
	PledgeID  int `json:"id"`
	UserID    string
	Item      Item      `json:"item"`
	Timestamp time.Time `json:"timestamp"`
}

type PledgeRequest struct {
	ItemId int `json:"itemId"`
}

type User struct {
	ID        string
	Country   string
	FirstName string
	Email     string
	AuthToken string
	Fullname  string
	IsAdmin   bool
}

// NewItem is a record of a (user) request to add a new item to the database of
// Items. 'Company' will be used if also adding in a new company to the
// database. IsPledge is true if the user also pledged the item
type NewItem struct {
	ID        int
	UserID    string
	IsPledge  bool
	Make      string
	Model     string
	Value     int
	Company   string
	CompanyID int
	Timestamp time.Time
}

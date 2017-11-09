package mdl

import (
	"time"
)

type ID string

type Company struct {
	Id   ID
	Name string
}

type Item struct {
	Id        ID
	Make      string
	Model     string
	CompanyID ID
}

// NewItem is a record of a (user) request to add a new item to the database of
// Items. 'Company' will be used if also adding in a new company to the
// database. IsPledge is true if the user also pledged the item
type NewItem struct {
	Id        ID
	UserID    ID
	IsPledge  bool
	Make      string
	Model     string
	Company   string
	CompanyID ID
	Timestamp time.Time
}

type Pledge struct {
	Id        ID
	UserID    ID
	ItemID    ID
	Timestamp time.Time
}

type ByTimeStamp []Pledge

func (b ByTimeStamp) Len() int           { return len(b) }
func (b ByTimeStamp) Less(i, j int) bool { return b[i].Timestamp.Before(b[j].Timestamp) }
func (b ByTimeStamp) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

type User struct {
	Id        ID
	Country   string
	FirstName string
	Email     string
	AuthToken string
	Fullname  string
}

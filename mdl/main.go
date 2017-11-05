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

type ItemUser struct {
	ItemID ID
	UserID ID
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

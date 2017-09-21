package mdl

import (
	"time"
)

type Id string

type CompanyRecord struct {
	Id   Id
	Name string
}

type ItemRecord struct {
	Id        Id
	Make      string
	Model     string
	CompanyId Id
}

type PledgeRecord struct {
	Id        Id
	UserId    Id
	ItemId    Id
	Timestamp time.Time
}

type ByTimeStamp []PledgeRecord

func (b ByTimeStamp) Len() int           { return len(b) }
func (b ByTimeStamp) Less(i, j int) bool { return b[i].Timestamp.Before(b[j].Timestamp) }
func (b ByTimeStamp) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

type UserRecord struct {
	Id        Id
	Country   string
	FirstName string
	Email     string
	AuthToken string
	Fullname  string
}

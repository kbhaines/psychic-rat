package mdl

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Id string

type CompanyRecord struct {
	Id   Id
	Name string
}

func NewCompany(id Id, name string) CompanyRecord {
	return CompanyRecord{Id: id, Name: name}
}

type ItemRecord struct {
	Id        Id
	Make      string
	Model     string
	CompanyId Id
}

func NewItem(make, model string, companyId Id) ItemRecord {
	return ItemRecord{Id: Id(uuid.NewV4().String()), Make: make, Model: model, CompanyId: companyId}
}

type PledgeRecord struct {
	Id        Id
	UserId    Id
	ItemId    Id
	Timestamp time.Time
}

func NewPledge(id Id, userId Id, itemId Id, timestamp time.Time) PledgeRecord {
	return PledgeRecord{Id: id, UserId: userId, ItemId: itemId, Timestamp: timestamp}
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

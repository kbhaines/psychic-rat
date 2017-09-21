package mdl

import (
	"fmt"
	"psychic-rat/mdl/company"
	"psychic-rat/mdl/item"
	"psychic-rat/mdl/user"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Id string

type CompanyRecord interface {
	Id() Id
	Name() string
}

type companyRecord struct {
	id   Id
	name string
}

func NewCompany(id Id, name string) CompanyRecord { return &companyRecord{id, name} }

func (r *companyRecord) Id() Id       { return r.id }
func (r *companyRecord) Name() string { return r.name }

type ItemRecord interface {
	Id() Id
	Make() string
	Model() string
	Company() company.Id
}

func NewItem(make string, model string, company company.Id) ItemRecord {
	id := Id(uuid.NewV4().String())
	return &itemRecord{id: id, make: make, model: model, company: company}
}

type itemRecord struct {
	id      Id
	make    string
	model   string
	company company.Id
}

func (r *itemRecord) Id() Id {
	return r.id
}

func (r *itemRecord) Make() string {
	return r.make
}

func (r *itemRecord) Model() string {
	return r.model
}

func (r *itemRecord) Company() company.Id {
	return r.company
}

func (r *itemRecord) String() string {
	return fmt.Sprintf("item: %v", *r)
}

type PledgeRecord interface {
	Id() Id
	UserId() user.Id
	ItemId() item.Id
	TimeStamp() time.Time
}

type ByTimeStamp []PledgeRecord

func (b ByTimeStamp) Len() int           { return len(b) }
func (b ByTimeStamp) Less(i, j int) bool { return b[i].TimeStamp().Before(b[j].TimeStamp()) }
func (b ByTimeStamp) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

type pledgeRecord struct {
	id        Id
	userId    user.Id
	itemId    item.Id
	timestamp time.Time
}

func NewPledge(userId user.Id, itemId item.Id, timestamp time.Time) PledgeRecord {
	return &pledgeRecord{id: Id(uuid.NewV4().String()), userId: userId, itemId: itemId, timestamp: timestamp}
}

func (r *pledgeRecord) Id() Id {
	return r.id
}

func (r *pledgeRecord) UserId() user.Id {
	return r.userId
}

func (r *pledgeRecord) ItemId() item.Id {
	return r.itemId
}

func (r *pledgeRecord) TimeStamp() time.Time {
	return r.timestamp
}

type UserRecord interface {
	Id() Id
	Country() string
	FirstName() string

	Email() string
	AuthToken() string
	Fullname() string
}

type userRecord struct {
	id        Id
	country   string
	firstName string
	email     string
	authToken string
	fullname  string
}

func NewUser(email string, country string, firstname string) UserRecord {
	return &userRecord{email: email, country: country, firstName: firstname}
}

func (r *userRecord) Id() Id {
	return r.id
}

func (r *userRecord) Country() string {
	return r.country
}

func (r *userRecord) FirstName() string {
	return r.firstName
}

func (r *userRecord) Email() string {
	return r.email
}

func (r *userRecord) AuthToken() string {
	return r.authToken
}

func (r *userRecord) Fullname() string {
	return r.fullname
}

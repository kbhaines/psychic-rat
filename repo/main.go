package repo

import (
	"psychic-rat/mdl/company"
	"psychic-rat/mdl/item"
	"psychic-rat/mdl/pledge"
	"psychic-rat/mdl/user"
)

type Companies interface {
	Create(company company.CompanyRecord) error
	GetCompanies() []company.CompanyRecord
	GetById(company.Id) (company.CompanyRecord, error)
}

type Items interface {
	Create(item item.ItemRecord) error
	GetById(id item.Id) (item.ItemRecord, error)
	GetAllByCompany(companyId company.Id) []item.ItemRecord
	Update(id item.Id, item item.ItemRecord)
	List() []item.ItemRecord
}

type Pledges interface {
	Create(pledge pledge.PledgeRecord) error
	GetById(id pledge.Id) (pledge.PledgeRecord, error)
	GetByUser(id user.Id) []pledge.Id
	List() []pledge.PledgeRecord
}

type Users interface {
	Create(user user.UserRecord) error
	GetById(id user.Id) (user.UserRecord, error)
}

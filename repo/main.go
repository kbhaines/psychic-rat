package repo

import (
	"psychic-rat/mdl/company"
	"psychic-rat/mdl/item"
	"psychic-rat/mdl/pledge"
	"psychic-rat/mdl/user"
)

type Companies interface {
	Create(company company.Record) error
	GetCompanies() []company.Record
	GetById(company.Id) (company.Record, error)
}

type Items interface {
	Create(item item.Record) error
	GetById(id item.Id) (item.Record, error)
	GetAllByCompany(companyId company.Id) []item.Record
	Update(id item.Id, item item.Record)
	List() []item.Record
}

type Pledges interface {
	Create(pledge pledge.Record) error
	GetById(id pledge.Id) (pledge.Record, error)
	GetByUser(id user.Id) []pledge.Id
	List() []pledge.Record
}

type Users interface {
	Create(user user.Record) error
	GetById(id user.Id) (user.Record, error)
}

package api

import "psychic-rat/mdl"

type Companies interface {
	Create(company mdl.CompanyRecord) error
	GetCompanies() []mdl.CompanyRecord
	GetById(mdl.Id) (*mdl.CompanyRecord, error)
}

type Items interface {
	Create(item mdl.ItemRecord) error
	GetById(id mdl.Id) (*mdl.ItemRecord, error)
	GetAllByCompany(companyId mdl.Id) []mdl.ItemRecord
	Update(id mdl.Id, item mdl.ItemRecord)
	List() []mdl.ItemRecord
}

type Pledges interface {
	Create(pledge mdl.PledgeRecord) error
	GetById(id mdl.Id) (*mdl.PledgeRecord, error)
	GetByUser(id mdl.Id) []mdl.Id
	List() []mdl.PledgeRecord
}

type Users interface {
	GetById(id mdl.Id) (*mdl.UserRecord, error)
	Create(user mdl.UserRecord) error
}

type Repos struct {
	Company Companies
	Item    Items
	Pledge  Pledges
	User    Users
}

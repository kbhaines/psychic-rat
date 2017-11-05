package api

import "psychic-rat/mdl"

type Companies interface {
	Create(company mdl.Company) error
	GetCompanies() []mdl.Company
	GetById(mdl.ID) (*mdl.Company, error)
}

type Items interface {
	Create(item mdl.Item) error
	GetById(id mdl.ID) (*mdl.Item, error)
	GetAllByCompany(companyId mdl.ID) []mdl.Item
	Update(id mdl.ID, item mdl.Item)
	List() []mdl.Item
}

type Pledges interface {
	Create(pledge mdl.Pledge) error
	GetById(id mdl.ID) (*mdl.Pledge, error)
	GetByUser(id mdl.ID) []mdl.ID
	List() []mdl.Pledge
}

type Users interface {
	GetById(id mdl.ID) (*mdl.User, error)
	Create(user mdl.User) error
}

type Repos struct {
	Company Companies
	Item    Items
	Pledge  Pledges
	User    Users
}

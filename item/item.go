package item

import (
	"psychic-rat/company"
)

type Id int

type Record struct {
	Id           Id
	Name         string
	Description  string
	Manufacturer company.Id
}

type Repo interface {
	Create(item Record) (Id, error)
	GetById(id Id) (Record, error)
	List() []Id
}


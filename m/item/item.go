package item

import (
	"psychic-rat/m/company"
	"github.com/satori/go.uuid"
)

type Id int

type Record interface {
	Id() Id
	Name() string
	Description() string
	Manufacturer() company.Id
}

type Repo interface {
	Create(item Record) (Id, error)
	GetById(id Id) (Record, error)
	List() []Id
}

func New(name string, description string, manufacturer company.Id) Record {
	return &record{id: Id(uuid.NewV4().String()), description: description, manufacturer: manufacturer}
}

type record struct {
	id           Id
	name         string
	description  string
	manufacturer company.Id
}

func (r *record) Id() Id {
	return r.id
}

func (r *record) Name() string {
	return r.name
}

func (r *record) Description() string {
	return r.description
}

func (r *record) Manufacturer() company.Id {
	return r.manufacturer
}

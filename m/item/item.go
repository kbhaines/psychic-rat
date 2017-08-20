package item

import (
	"psychic-rat/m/company"
	"github.com/satori/go.uuid"
)

type Id string

type Record interface {
	Id() Id
	Make() string
	Model() string
	Manufacturer() company.Id
}

type Repo interface {
	Create(item Record) (Id, error)
	GetById(id Id) (Record, error)
	List() []Id
}

func New(make string, model string, manufacturer company.Id) Record {
	return &record{id: Id(uuid.NewV4().String()), model: model, manufacturer: manufacturer}
}

type record struct {
	id           Id
	make         string
	model        string
	manufacturer company.Id
}

func (r *record) Id() Id {
	return r.id
}

func (r *record) Make() string {
	return r.make
}

func (r *record) Model() string {
	return r.model
}

func (r *record) Manufacturer() company.Id {
	return r.manufacturer
}

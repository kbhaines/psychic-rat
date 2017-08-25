package item

import (
	"psychic-rat/mdl/company"
	"fmt"
	"github.com/satori/go.uuid"
)

type Id string

type Record interface {
	Id() Id
	Make() string
	Model() string
	Company() company.Id
}

func New(make string, model string, company company.Id) Record {
	id := Id(uuid.NewV4().String())
	return &record{id: id, make: make, model: model, company: company}
}

type record struct {
	id      Id
	make    string
	model   string
	company company.Id
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

func (r *record) Company() company.Id {
	return r.company
}

func (r *record) String() string {
	return fmt.Sprintf("item: %v", *r)
}

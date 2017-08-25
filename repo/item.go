package repo

import "psychic-rat/mdl/company"
import (
	"psychic-rat/mdl/item"
)

type Items interface {
	Create(item item.Record) error
	GetById(id item.Id) (item.Record, error)
	GetAllByCompany(companyId company.Id) []item.Record
	Update(id item.Id, item item.Record)
	List() []item.Id
}

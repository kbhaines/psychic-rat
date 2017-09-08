package repo

import "psychic-rat/mdl/company"

type Companies interface {
	Create(company company.Record) error
	GetCompanies() []company.Record
	GetById(company.Id) (company.Record, error)
}

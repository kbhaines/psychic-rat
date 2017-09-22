package repo

import (
	"fmt"
	"psychic-rat/mdl"
)

var companyRepo Companies = &companyRepoMap{make(map[mdl.Id]mdl.CompanyRecord)}

func getCompanyRepoMapImpl() Companies {
	return companyRepo
}

type companyRepoMap struct {
	records map[mdl.Id]mdl.CompanyRecord
}

func (r *companyRepoMap) Create(i mdl.CompanyRecord) error {
	if _, found := r.records[i.Id]; found {
		return fmt.Errorf("company id %v exists", i.Id)
	}
	r.records[i.Id] = i
	return nil
}

func (r *companyRepoMap) GetCompanies() (companies []mdl.CompanyRecord) {
	for _, c := range r.records {
		companies = append(companies, c)
	}
	return companies
}

func (r *companyRepoMap) GetById(id mdl.Id) (*mdl.CompanyRecord, error) {
	if rec, exists := r.records[id]; exists {
		return &rec, nil
	}
	return nil, fmt.Errorf("company id %v not found", id)
}

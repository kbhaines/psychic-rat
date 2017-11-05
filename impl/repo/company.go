package repo

import (
	"fmt"
	"psychic-rat/mdl"
)

var companyRepo = &companyRepoMap{make(map[mdl.ID]mdl.Company)}

type companyRepoMap struct {
	records map[mdl.ID]mdl.Company
}

func (r *companyRepoMap) Create(i mdl.Company) error {
	if _, found := r.records[i.Id]; found {
		return fmt.Errorf("company id %v exists", i.Id)
	}
	r.records[i.Id] = i
	return nil
}

func (r *companyRepoMap) GetCompanies() (companies []mdl.Company) {
	for _, c := range r.records {
		companies = append(companies, c)
	}
	return companies
}

func (r *companyRepoMap) GetById(id mdl.ID) (*mdl.Company, error) {
	if rec, exists := r.records[id]; exists {
		return &rec, nil
	}
	return nil, fmt.Errorf("company id %v not found", id)
}

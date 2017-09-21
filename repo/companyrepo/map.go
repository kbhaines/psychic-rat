package companyrepo

import (
	"fmt"
	"psychic-rat/mdl"
	"psychic-rat/repo"
)

// declare that we implement Repo interface
var companyRepo repo.Companies = &repoMap{make(map[mdl.Id]mdl.CompanyRecord)}

func GetCompanyRepoMapImpl() repo.Companies {
	return companyRepo
}

type repoMap struct {
	records map[mdl.Id]mdl.CompanyRecord
}

func (r *repoMap) Create(i mdl.CompanyRecord) error {
	if _, found := r.records[i.Id]; found {
		return fmt.Errorf("company id %v exists", i.Id)
	}
	r.records[i.Id] = i
	return nil
}

func (r *repoMap) GetCompanies() (companies []mdl.CompanyRecord) {
	for _, c := range r.records {
		companies = append(companies, c)
	}
	return companies
}

func (r *repoMap) GetById(id mdl.Id) (*mdl.CompanyRecord, error) {
	if rec, exists := r.records[id]; exists {
		return &rec, nil
	}
	return nil, fmt.Errorf("company id %v not found", id)
}

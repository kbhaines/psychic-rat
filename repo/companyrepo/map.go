package companyrepo

import (
	"psychic-rat/mdl/company"
	"psychic-rat/repo"
	"fmt"
)

// declare that we implement Repo interface
var companyRepo repo.Companies = &repoMap{make(map[company.Id]company.Record)}

func GetCompanyRepoMapImpl() repo.Companies {
	return companyRepo
}

type repoMap struct {
	records map[company.Id]company.Record
}

func (r *repoMap) Create(i company.Record) error {
	if _, found := r.records[i.Id()]; found {
		return fmt.Errorf("company id %v exists", i.Id())
	}
	r.records[i.Id()] = i
	return nil
}

func (r *repoMap) GetCompanies() (companies []company.Record) {
	for _, c := range r.records {
		companies = append(companies, c)
	}
	return companies
}

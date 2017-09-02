package ctr

import (
	"psychic-rat/mdl/company"
	"psychic-rat/repo/companyrepo"
)

type CompanyController interface {
	GetCompanies() []company.Record
}

type companyController struct{}

var companyRepo = companyrepo.GetCompanyRepoMapImpl()

func (c *companyController) GetCompanies() []company.Record {
	return companyRepo.GetCompanies()
}

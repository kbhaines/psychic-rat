package main

import (
	"net/http"
	"psychic-rat/api"
	"psychic-rat/repo/companyrepo"
	"psychic-rat/mdl/company"
)


func main() {
	repo := companyrepo.GetCompanyRepoMapImpl()
	repo.Create(company.New(company.Id("1"), "bigco1"))
	repo.Create(company.New(company.Id("2"), "bigco2"))
	repo.Create(company.New(company.Id("3"), "bigco3"))
	http.HandleFunc("/api/v1/company", api.CompanyHandler)
	http.ListenAndServe("localhost:8080", nil)
}

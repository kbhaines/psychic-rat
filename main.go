package main

import (
	"net/http"
	"psychic-rat/api"
	"psychic-rat/repo/companyrepo"
	"psychic-rat/mdl/company"
	"psychic-rat/repo/itemrepo"
	"psychic-rat/mdl/item"
)


func main() {
	companies := companyrepo.GetCompanyRepoMapImpl()
	companies.Create(company.New(company.Id("1"), "bigco1"))
	companies.Create(company.New(company.Id("2"), "bigco2"))
	companies.Create(company.New(company.Id("3"), "bigco3"))

	items := itemrepo.GetItemRepoMapImpl()
	items.Create(item.New("phone", "abc", company.Id("1")))
	items.Create(item.New("phone", "xyz", company.Id("1")))
	items.Create(item.New("tablet", "gt1", company.Id("1")))
	items.Create(item.New("tablet", "tab4", company.Id("2")))
	items.Create(item.New("tablet", "tab8", company.Id("2")))

	http.HandleFunc("/api/v1/company", api.CompanyHandler)
	http.HandleFunc("/api/v1/item", api.ItemHandler)
	http.HandleFunc("/api/v1/pledge", api.PledgeHandler)
	
	http.ListenAndServe("localhost:8080", nil)
}

package repo

import (
	"psychic-rat/mdl"
	"psychic-rat/repo"
)

func GetCompanyRepoMapImpl() repo.Companies {
	return companyRepo
}

func GetItemRepoMapImpl() repo.Items {
	return itemRepo
}

func GetPledgeRepoMapImpl() repo.Pledges {
	return pledgeRepo
}

func GetUserRepoMapImpl() repo.Users {
	return userRepo
}

func GetRepos() repo.Repos {
	return repo.Repos{
		Company: companyRepo,
		Item:    itemRepo,
		Pledge:  pledgeRepo,
		User:    userRepo,
	}
}

func init() {
	setupMockData()
}

func setupMockData() {
	companies := GetRepos().Company
	companies.Create(mdl.NewCompany(mdl.Id("1"), "bigco1"))
	companies.Create(mdl.NewCompany(mdl.Id("2"), "bigco2"))
	companies.Create(mdl.NewCompany(mdl.Id("3"), "bigco3"))

	items := GetRepos().Item
	items.Create(mdl.NewItem("1", "phone", "abc", mdl.Id("1")))
	items.Create(mdl.NewItem("2", "phone", "xyz", mdl.Id("1")))
	items.Create(mdl.NewItem("3", "tablet", "gt1", mdl.Id("1")))
	items.Create(mdl.NewItem("4", "tablet", "tab4", mdl.Id("2")))
	items.Create(mdl.NewItem("5", "tablet", "tab8", mdl.Id("2")))

	users := GetRepos().User
	users.Create(mdl.UserRecord{Id: mdl.Id("testuser1"), Email: "testuser1@gmail.com", FirstName: "Kevin"})
}

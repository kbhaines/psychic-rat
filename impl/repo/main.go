package repo

import (
	"psychic-rat/mdl"
)

func GetCompanyRepoMapImpl() *companyRepoMap {
	return companyRepo
}

func GetItemRepoMapImpl() *itemRepoMap {
	return itemRepo
}

func GetPledgeRepoMapImpl() *pledgeRepoMap {
	return pledgeRepo
}

func GetUserRepoMapImpl() *userRepoMap {
	return userRepo
}

func init() {
	setupMockData()
}

func setupMockData() {
	companyRepo.Create(mdl.NewCompany(mdl.Id("1"), "bigco1"))
	companyRepo.Create(mdl.NewCompany(mdl.Id("2"), "bigco2"))
	companyRepo.Create(mdl.NewCompany(mdl.Id("3"), "bigco3"))

	itemRepo.Create(mdl.NewItem("1", "phone", "abc", mdl.Id("1")))
	itemRepo.Create(mdl.NewItem("2", "phone", "xyz", mdl.Id("1")))
	itemRepo.Create(mdl.NewItem("3", "tablet", "gt1", mdl.Id("1")))
	itemRepo.Create(mdl.NewItem("4", "tablet", "tab4", mdl.Id("2")))
	itemRepo.Create(mdl.NewItem("5", "tablet", "tab8", mdl.Id("2")))

	userRepo.Create(mdl.UserRecord{Id: mdl.Id("testuser1"), Email: "testuser1@gmail.com", Fullname: "Kevin"})
	userRepo.Create(mdl.UserRecord{Id: mdl.Id("testuser2"), Email: "testuser2@gmail.com", Fullname: "Brian"})
	userRepo.Create(mdl.UserRecord{Id: mdl.Id("testuser3"), Email: "testuser3@gmail.com", Fullname: "Lynnette"})
}

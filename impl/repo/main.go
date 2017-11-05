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

	companyRepo.Create(mdl.Company{mdl.ID("1"), "bigco1"})
	companyRepo.Create(mdl.Company{mdl.ID("2"), "bigco2"})
	companyRepo.Create(mdl.Company{mdl.ID("3"), "bigco3"})

	itemRepo.Create(mdl.Item{Id: mdl.ID("1"), Make: "phone", Model: "abc", CompanyID: mdl.ID("1")})
	itemRepo.Create(mdl.Item{Id: mdl.ID("2"), Make: "phone", Model: "xyz", CompanyID: mdl.ID("1")})
	itemRepo.Create(mdl.Item{Id: mdl.ID("3"), Make: "tablet", Model: "gt1", CompanyID: mdl.ID("1")})
	itemRepo.Create(mdl.Item{Id: mdl.ID("4"), Make: "tablet", Model: "tab4", CompanyID: mdl.ID("2")})
	itemRepo.Create(mdl.Item{Id: mdl.ID("5"), Make: "tablet", Model: "tab8", CompanyID: mdl.ID("2")})

	userRepo.Create(mdl.User{Id: mdl.ID("testuser1"), Email: "testuser1@gmail.com", Fullname: "Kevin"})
	userRepo.Create(mdl.User{Id: mdl.ID("testuser2"), Email: "testuser2@gmail.com", Fullname: "Brian"})
	userRepo.Create(mdl.User{Id: mdl.ID("testuser3"), Email: "testuser3@gmail.com", Fullname: "Lynnette"})
}

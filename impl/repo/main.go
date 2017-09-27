package repo

import "psychic-rat/repo"

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

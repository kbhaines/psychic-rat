package api

import (
	a "psychic-rat/api"
	"psychic-rat/repo"
)

func GetApis(repos repo.Repos) a.Api {
	return a.Api{
		Company: getRepoCompanyApi(repos),
		Item:    getRepoItemApi(repos),
		Pledge:  getRepoPledgeApi(repos),
	}
}

func getRepoCompanyApi(repos repo.Repos) a.CompanyApi {
	return &companyApiRepoImpl{repos: repos}
}

func getRepoItemApi(repos repo.Repos) a.ItemApi {
	return &itemRepoApi{repos: repos}
}

func getRepoPledgeApi(repos repo.Repos) a.PledgeApi {
	return &pledgeApiRepoImpl{repos: repos}
}

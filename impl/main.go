package impl

import a "psychic-rat/api"
import "psychic-rat/repo"
import "psychic-rat/impl/repo"

func init() {
	repos = repo.Repos{
		Company: GetCompanyRepoMapImpl(),
		Item:    GetItemRepoMapImpl(),
		Pledge:  GetPledgeRepoMapImpl(),
		User:    GetUserRepoMapImpl(),
	}

	api = a.Api{
		Company: GetRepoCompanyApi(repos),
		Item:    GetRepoItemApi(repos),
		Pledge:  GetRepoPledgeApi(repos),
	}

}

var api a.Api
var repos repo.Repos

func Get() a.Api {
	return api
}

func GetRepos() repo.Repos {
	return repos
}

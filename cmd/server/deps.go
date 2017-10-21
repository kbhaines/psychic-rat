package main

import "psychic-rat/impl/repo"
import "psychic-rat/impl/api"

func init() {
	repos := api.Repos{
		Company: repo.GetCompanyRepoMapImpl(),
		Item:    repo.GetItemRepoMapImpl(),
		Pledge:  repo.GetPledgeRepoMapImpl(),
		User:    repo.GetUserRepoMapImpl(),
	}
	apis = Api{
		Company: api.NewCompanyApi(repos),
		Item:    api.NewItemApi(repos),
		Pledge:  api.NewPledgeApi(repos),
		User:    api.NewUserApi(repos),
	}
}

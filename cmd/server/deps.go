package main

import "psychic-rat/impl/repo"
import "psychic-rat/impl/api"

func init() {
	repos := repo.GetRepos()
	apis = Api{
		Company: api.NewCompanyApi(repos),
		Item:    api.NewItemApi(repos),
		Pledge:  api.NewPledgeApi(repos),
		User:    api.NewUserApi(repos),
	}
}

package web

import (
	"psychic-rat/auth0"
	"psychic-rat/impl/api"
	"psychic-rat/impl/repo"
	"psychic-rat/sqldb"
)

func init() {
	repos := api.Repos{
		Company: repo.GetCompanyRepoMapImpl(),
		Item:    repo.GetItemRepoMapImpl(),
		Pledge:  repo.GetPledgeRepoMapImpl(),
		User:    repo.GetUserRepoMapImpl(),
	}
	apis = API{
		Company: api.NewCompanyApi(repos),
		Item:    api.NewItemApi(repos),
		Pledge:  api.NewPledgeApi(repos),
		User:    api.NewUserApi(repos),
	}
	auth0.Init(apis.User)

	var err error
	db, err = sqldb.NewDB("pr.dat")
	if err != nil {
		panic("unable to init db: " + err.Error())
	}

}

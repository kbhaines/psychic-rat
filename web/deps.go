package web

import (
	"psychic-rat/auth0"
	"psychic-rat/sqldb"
)

func init() {
	var err error
	db, err = sqldb.NewDB("pr.dat")
	if err != nil {
		panic("unable to init db: " + err.Error())
	}

	apis = API{
		Company: db,
		Item:    db,
		Pledge:  db,
		User:    db,
	}
	auth0.Init(apis.User)
}

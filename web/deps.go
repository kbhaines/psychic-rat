package web

import (
	"psychic-rat/auth0"
	"psychic-rat/sqldb"
)

func init() {
	// TODO: smells a bit, as it gets used in tests by default

	var err error
	db, err = sqldb.NewDB("pr.dat")
	if err != nil {
		panic("unable to init db: " + err.Error())
	}

	apis = API{
		Company: db,
		Item:    db,
		NewItem: db,
		Pledge:  db,
		User:    db,
	}
	auth0.Init(apis.User)
}

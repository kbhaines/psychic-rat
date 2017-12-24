package web

import (
	"psychic-rat/auth0"
	"psychic-rat/sqldb"
	"psychic-rat/web/admin"
)

func init() {
	// TODO: smells a bit, as it gets used in tests by default

	db, err := sqldb.OpenDB("pr.dat")
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
	admin.InitDeps(db, db, db, db)
}

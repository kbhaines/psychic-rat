package web

import (
	"psychic-rat/auth0"
	"psychic-rat/sqldb"
)

func init() {
	apis = API{}
	auth0.Init(apis.User)

	var err error
	db, err = sqldb.NewDB("pr.dat")
	if err != nil {
		panic("unable to init db: " + err.Error())
	}

}

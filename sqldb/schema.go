package sqldb

import (
	"database/sql"
	"fmt"
	"log"
)

const schemaVersion = 8

var updateFuncs = map[int]func(*sql.DB) error{
	6: update6,
	8: update8,
}

func createSchema(db *sql.DB) error {
	stmt := `
	create table companies (id integer primary key, name string);

	create table items (id integer primary key,
	  make string, 
	  model string, 
	  companyID integer,
      usdValue integer,
  	  newItemID integer);
	
	create view itemsCompany as select i.*, c.name 'companyName' from items i, companies c where i.companyID = c.id;

	create table currencies (
		id integer primary key,
		ident string,
		usdConversion float);

	create table newItems(id integer primary key, 
	  userID integer,
	  isPledge boolean,
	  make string, 
	  model string, 
	  company string,
	  companyID integer,
	  currencyId integer,
	  currencyValue integer,
	  timestamp integer,
 	  used boolean );

	create table users (id string primary key,
	  fullName string,
	  firstName string,
	  country string,
	  email string,
  	  isAdmin bool);

	create table pledges (id integer primary key,
	  userID integer,
	  itemID integer,
	  usdValue integer,
	  timestamp integer);

	create view userPledges as select p.id pledgeID, p.userID, p.timestamp, 
		i.id itemID, i.Make, i.Model, c.id companyID, c.name 
		from pledges p, items i, companies c 
		where p.itemID = i.id and i.companyID = c.id;

	`
	_, err := db.Exec(stmt)
	if err != nil {
		log.Fatalf("unable to create schema: %v", err)
	}
	setVersion(db, schemaVersion)
	return nil
}

func schemaUpdates(db *sql.DB) {
	version := getVersion(db)
	for v := version + 1; v <= schemaVersion; v++ {
		updateFunc, ok := updateFuncs[v]
		if ok {
			err := updateFunc(db)
			if err != nil {
				log.Fatal(err)
			}
			setVersion(db, v)
		}
	}
	return
}

func getVersion(db *sql.DB) int {
	row := db.QueryRow("pragma user_version")
	var version int
	err := row.Scan(&version)
	if err != nil {
		log.Fatalf("unable to determine version: %v", err)
	}
	return version
}

func setVersion(db *sql.DB, v int) {
	_, err := db.Exec(fmt.Sprintf("pragma user_version = %d", v))
	if err != nil {
		log.Fatalf("unable to set schema version: %v", err)
	}
}

func update6(db *sql.DB) error {
	log.Print("updating schema to version 6")
	return nil
}

func update8(db *sql.DB) error {
	log.Print("updating schema to version 8")
	return nil
}

package sqldb

import (
	"database/sql"
	"fmt"
	"log"
)

const schemaVersion = 1

func createSchema(db *sql.DB) error {
	log.Print("creating initial schema")
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
	vers := getVersion(db)
	if vers < 1 {
		createSchema(db)
	}
	if vers < 2 {
		updateV2(db)
	}
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

func updateV2(db *sql.DB) {
	log.Print("updating schema to version 2")
	setVersion(db, 2)
	return
}

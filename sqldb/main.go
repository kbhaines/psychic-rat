package sqldb

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DB *sql.DB

func NewDB(name string) (DB, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}
	return DB(db), nil
}

package sqldb

import (
	"database/sql"
	"os"
	"psychic-rat/mdl"
	"psychic-rat/types"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func NewDB(name string) (*DB, error) {
	os.Remove(name)
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}
	createSchema(db)
	return &DB{db}, nil
}

func createSchema(db *sql.DB) error {
	stmt := `
	create table companies (id integer primary key, name string);
	`
	_, err := db.Exec(stmt)
	return err
}

func (d *DB) NewCompany(c mdl.Company) error {
	stmt, err := d.Prepare("insert into companies(name) values(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.Name)
	return err
}

func (d *DB) GetCompanies() (types.CompanyListing, error) {
	result := types.CompanyListing{}
	rows, err := d.Query("select id, name from companies")
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return result, err
		}
		result.Companies = append(result.Companies, types.Company{Id: mdl.ID(strconv.Itoa(id)), Name: name})
	}
	return result, nil
}

func (d *DB) GetById(mdl.ID) (types.Company, error) {
	panic("not implemented")
}

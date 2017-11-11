package sqldb

import (
	"database/sql"
	"os"
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
	err = createSchema(db)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func createSchema(db *sql.DB) error {
	stmt := `
	create table companies (id integer primary key, name string);
	create table items (id integer primary key,
	  make string, model string, companyId integer);
	
	create table newItems(id integer primary key, 
	  userId integer,
	  isPledge boolean,
	  make string, 
	  model string, 
	  company string,
	  companyId integer,
	  timestamp date
	  );
	`
	_, err := db.Exec(stmt)
	return err
}

func (d *DB) NewCompany(c types.Company) error {
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
		var co types.Company
		err = rows.Scan(&co.Id, &co.Name)
		if err != nil {
			return result, err
		}
		result.Companies = append(result.Companies, co)
	}
	return result, nil
}

func (d *DB) GetCompany(id int) (types.Company, error) {
	result := types.Company{}
	err := d.QueryRow("select id, name from companies where id = "+strconv.Itoa(id)).Scan(&result.Id, &result.Name)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (d *DB) ListItems() (types.ItemReport, error) {
	ir := types.ItemReport{}
	rows, err := d.Query("select id, make, model, companyId from items")
	if err != nil {
		return ir, err
	}
	for rows.Next() {
		item := types.Item{}
		var coid int
		err = rows.Scan(&item.Id, &item.Make, &item.Model, &coid)
		if err != nil {
			return ir, err
		}
		ir.Items = append(ir.Items, item)
	}
	return ir, nil
}

func (d *DB) GetItem(id int) (types.Item, error) {
	panic("not implemented")
}

func (d *DB) AddNewItem(i types.NewItem) error {
	s, err := d.Prepare("insert into newItems(userId, isPledge, make, model, company, companyId) values (?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer s.Close()
	_, err = s.Exec(i.UserID, i.IsPledge, i.Make, i.Model, i.Company, i.CompanyID)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) ListNewItems() ([]types.NewItem, error) {
	panic("not implemented")
}

func (d *DB) ApproveItem(id int) error {
	panic("not implemented")
}

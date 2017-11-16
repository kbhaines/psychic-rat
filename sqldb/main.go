package sqldb

import (
	"database/sql"
	"fmt"
	"os"
	"psychic-rat/mdl"
	"psychic-rat/types"
	"strconv"
	"time"

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
	  make string, 
	  model string, 
	  companyId integer);
	
	create table newItems(id integer primary key, 
	  userId integer,
	  isPledge boolean,
	  make string, 
	  model string, 
	  company string,
	  companyId integer,
	  timestamp integer);

	create table users (id string primary key,
	  fullName string,
	  firstName string,
	  country string,
	  email string);

	create table pledges (id integer primary key,
	  userId integer,
	  itemId integer,
	  timestamp integer);
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

func (d *DB) GetCompanies() ([]types.Company, error) {
	result := []types.Company{}
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
		result = append(result, co)
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

func (d *DB) ListItems() ([]types.Item, error) {
	ir := []types.Item{}
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
		ir = append(ir, item)
	}
	return ir, nil
}

func (d *DB) GetItem(id int) (types.Item, error) {
	i := types.Item{}
	var companyID int
	err := d.QueryRow("select id, make, model, companyId from items where id = ?", id).Scan(&i.Id, &i.Make, &i.Model, &companyID)
	if err != nil {
		return i, fmt.Errorf("could not get item %d: %v ", id, err)
	}
	return i, nil
}

func (d *DB) AddItem(i types.Item) (*types.Item, error) {
	s, err := d.Prepare("insert into items(make, model, companyId) values (?,?,?)")
	if err != nil {
		return nil, err
	}
	defer s.Close()
	r, err := s.Exec(i.Make, i.Model, i.Company.Id)
	if err != nil {
		return nil, err
	}
	lastId, err := r.LastInsertId()
	if err != nil {
		return nil, err
	}
	new := i
	new.Id = int(lastId)
	return &new, nil
}

func (d *DB) AddNewItem(i types.NewItem) (*types.NewItem, error) {
	s, err := d.Prepare("insert into newItems(userId, isPledge, make, model, company, companyId, timestamp) values (?,?,?,?,?,?,?)")
	if err != nil {
		return nil, err
	}
	defer s.Close()
	timestamp := time.Now().Truncate(time.Second)
	r, err := s.Exec(i.UserID, i.IsPledge, i.Make, i.Model, i.Company, i.CompanyID, timestamp.Unix())
	if err != nil {
		return nil, err
	}
	lastId, err := r.LastInsertId()
	if err != nil {
		return nil, err
	}
	new := i
	new.Timestamp = timestamp
	new.Id = int(lastId)
	return &new, nil
}

func (d *DB) ListNewItems() ([]types.NewItem, error) {
	result := []types.NewItem{}
	rows, err := d.Query("select id, userId, isPledge, make, model, company, companyId, timestamp from newItems")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var n types.NewItem
		var timestamp int64
		err = rows.Scan(&n.Id, &n.UserID, &n.IsPledge, &n.Make, &n.Model, &n.Company, &n.CompanyID, &timestamp)
		if err != nil {
			return result, err
		}
		n.Timestamp = time.Unix(timestamp, 0)
		result = append(result, n)
	}
	return result, nil
}

func (d *DB) ApproveItem(id int) error {
	panic("not implemented")
}

func (d *DB) GetUser(userId string) (*mdl.User, error) {
	result := mdl.User{}
	err := d.QueryRow("select id, fullname, firstName, country, email from users where id= \""+userId+"\"").Scan(&result.Id, &result.Fullname, &result.FirstName, &result.Country, &result.Email)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (d *DB) CreateUser(u mdl.User) error {
	stmt, err := d.Prepare("insert into users(id, fullName, firstName, country, email) values (?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(u.Id, u.Fullname, u.FirstName, u.Country, u.Email)
	return nil
}

func (d *DB) NewPledge(itemId int, userId string) (int, error) {
	timestamp := time.Now().Truncate(time.Second)
	stmt, err := d.Prepare("insert into pledges(itemId, userId, timestamp) values (?,?,?)")
	if err != nil {
		return 0, err
	}
	r, err := stmt.Exec(itemId, userId, timestamp)
	if err != nil {
		return 0, fmt.Errorf("could not insert pledge for item %d for user %d: %v", itemId, userId, err)
	}
	lastId, err := r.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("no id returned for pledge for item %d for user %d: %v", itemId, userId, err)
	}
	return int(lastId), nil
}

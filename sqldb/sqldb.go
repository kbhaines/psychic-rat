package sqldb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"psychic-rat/types"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
	insertUser   *sql.Stmt
	insertPledge *sql.Stmt
}

func NewDB(name string) (*DB, error) {
	if _, err := os.Stat(name); os.IsExist(err) {
		panic(fmt.Sprintf("refusing to create when DB %s already exists", name))
	}
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}
	err = createSchema(db)
	if err != nil {
		return nil, err
	}
	return setupDB(db)
}

func OpenDB(name string) (*DB, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}
	return setupDB(db)
}

func setupDB(db *sql.DB) (*DB, error) {

	// TODO: insert schema version checks here

	_, err := db.Exec("PRAGMA synchronous = OFF")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		return nil, err
	}

	insertUser, err := db.Prepare("insert into users(id, fullName, firstName, country, email, isAdmin) values (?,?,?,?,?,?)")
	insertPledge, err := db.Prepare("insert into pledges(itemID, userID, timestamp) values (?,?,?)")
	if err != nil {
		return nil, err
	}

	return &DB{db, insertUser, insertPledge}, nil
}

func createSchema(db *sql.DB) error {
	stmt := `
	create table companies (id integer primary key, name string);

	create table items (id integer primary key,
	  make string, 
	  model string, 
	  companyID integer,
	  currencyID integer,
	  value integer,
      usdValue integer );
	
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
	  timestamp integer);

	create table users (id string primary key,
	  fullName string,
	  firstName string,
	  country string,
	  email string,
  	  isAdmin bool);

	create table pledges (id integer primary key,
	  userID integer,
	  itemID integer,
	  timestamp integer);

	create view userPledges as select p.id pledgeID, p.userID, p.timestamp, 
		i.id itemID, i.Make, i.Model, c.id companyID, c.name 
		from pledges p, items i, companies c 
		where p.itemID = i.id and i.companyID = c.id;
	`
	_, err := db.Exec(stmt)
	return err
}

func (d *DB) AddCompany(c types.Company) (*types.Company, error) {
	r, err := d.Exec("insert into companies(name) values(?)", c.Name)
	if err != nil {
		return nil, err
	}
	lastID, err := r.LastInsertId()
	if err != nil {
		return nil, err
	}
	c.ID = int(lastID)
	return &c, nil
}

func (d *DB) ListCompanies() ([]types.Company, error) {
	result := []types.Company{}
	rows, err := d.Query("select id, name from companies")
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var co types.Company
		err = rows.Scan(&co.ID, &co.Name)
		if err != nil {
			return result, err
		}
		result = append(result, co)
	}
	return result, nil
}

func (d *DB) GetCompany(id int) (types.Company, error) {
	result := types.Company{}
	err := d.QueryRow("select id, name from companies where id = ?", id).Scan(&result.ID, &result.Name)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (d *DB) ListItems() ([]types.Item, error) {
	ir := []types.Item{}
	rows, err := d.Query("select id, make, model, companyID, companyName, currencyID, value, usdValue from itemsCompany order by companyName, make, model")
	if err != nil {
		return ir, err
	}
	for rows.Next() {
		item := types.Item{}
		err = rows.Scan(&item.ID, &item.Make, &item.Model, &item.Company.ID, &item.Company.Name, &item.CurrencyID, &item.Value, &item.USDValue)
		if err != nil {
			return ir, err
		}
		ir = append(ir, item)
	}
	return ir, nil
}

func (d *DB) AddCurrency(c types.Currency) (*types.Currency, error) {
	r, err := d.Exec("insert into currencies(ident, usdConversion) values (?,?)", c.Ident, c.ConversionToUSD)
	if err != nil {
		return nil, err
	}
	lastID, err := r.LastInsertId()
	if err != nil {
		return nil, err
	}
	c.ID = int(lastID)
	return &c, nil
}

func (d *DB) ListCurrencies() ([]types.Currency, error) {
	cs := []types.Currency{}
	rows, err := d.Query("select id, ident, usdConversion from currencies order by id")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		c := types.Currency{}
		err = rows.Scan(&c.ID, &c.Ident, &c.ConversionToUSD)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}

func (d *DB) getCurrency(id int) (*types.Currency, error) {
	c := types.Currency{}
	err := d.QueryRow("select id, ident, usdConversion from currencies where id = ?", id).Scan(&c.ID, &c.Ident, &c.ConversionToUSD)
	if err != nil {
		return nil, fmt.Errorf("could not get currency %d", id)
	}
	return &c, nil
}

func (d *DB) GetItem(id int) (types.Item, error) {
	i := types.Item{}
	err := d.QueryRow("select id, make, model, companyID, companyName, currencyID, value, usdValue from itemsCompany where id = ?", id).Scan(&i.ID, &i.Make, &i.Model, &i.Company.ID, &i.Company.Name, &i.CurrencyID, &i.Value, &i.USDValue)
	if err != nil {
		return i, fmt.Errorf("could not get item %d: %v ", id, err)
	}
	return i, nil
}

func (d *DB) AddItem(i types.Item) (*types.Item, error) {
	c, err := d.getCurrency(i.CurrencyID)
	if err != nil {
		return nil, fmt.Errorf("could not get currency for new item %v: %v", i, err)
	}
	i.USDValue = int(c.ConversionToUSD * float64(i.Value))

	log.Printf("i,c = %+v %+v\n", i, c)
	r, err := d.Exec("insert into items(make, model, companyID, currencyID, value, usdValue) values (?,?,?,?,?,?)", i.Make, i.Model, i.Company.ID, i.CurrencyID, i.Value, i.USDValue)
	if err != nil {
		return nil, err
	}
	lastID, err := r.LastInsertId()
	if err != nil {
		return nil, err
	}
	new := i
	new.ID = int(lastID)
	return &new, nil
}

func (d *DB) AddNewItem(i types.NewItem) (*types.NewItem, error) {
	timestamp := time.Now().Truncate(time.Second)
	_, err := d.getCurrency(i.CurrencyID)
	if err != nil {
		return nil, err
	}
	r, err := d.Exec("insert into newItems(userID, isPledge, make, model, company, companyID, currencyID, currencyValue, timestamp) values (?,?,?,?,?,?,?,?,?)", i.UserID, i.IsPledge, i.Make, i.Model, i.Company, i.CompanyID, i.CurrencyID, i.Value, timestamp.Unix())
	if err != nil {
		return nil, err
	}
	lastID, err := r.LastInsertId()
	if err != nil {
		return nil, err
	}
	new := i
	new.Timestamp = timestamp
	new.ID = int(lastID)
	return &new, nil
}

func (d *DB) ListNewItems() ([]types.NewItem, error) {
	result := []types.NewItem{}
	rows, err := d.Query("select id, userID, isPledge, make, model, company, companyID, currencyID, currencyValue, timestamp from newItems")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var n types.NewItem
		var timestamp int64
		err = rows.Scan(&n.ID, &n.UserID, &n.IsPledge, &n.Make, &n.Model, &n.Company, &n.CompanyID, &n.CurrencyID, &n.Value, &timestamp)
		if err != nil {
			return result, err
		}
		n.Timestamp = time.Unix(timestamp, 0)
		result = append(result, n)
	}
	return result, nil
}

func (d *DB) DeleteNewItem(id int) error {
	_, err := d.Exec("delete from newItems where id=?", id)
	return err
}

func (d *DB) getNewItem(id int) (*types.NewItem, error) {
	i := types.NewItem{}
	var timestamp int64
	err := d.QueryRow("select id, userID, isPledge, make, model, company, companyID, timestamp from newItems where id = ?", id).Scan(&i.ID,
		&i.UserID, &i.IsPledge, &i.Make, &i.Model, &i.Company, &i.CompanyID, &timestamp)
	if err != nil {
		return nil, err
	}
	i.Timestamp = time.Unix(timestamp, 0)
	return &i, nil
}

func (d *DB) GetUser(userID string) (*types.User, error) {
	u := types.User{}
	err := d.QueryRow("select id, fullname, firstName, country, email, isAdmin from users where id=?", userID).Scan(&u.ID, &u.Fullname, &u.FirstName, &u.Country, &u.Email, &u.IsAdmin)
	if err != nil {
		return &u, err
	}
	return &u, nil
}

func (d *DB) AddUser(u types.User) error {
	_, err := d.insertUser.Exec(u.ID, u.Fullname, u.FirstName, u.Country, u.Email, u.IsAdmin)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) AddPledge(itemID int, userID string) (*types.Pledge, error) {
	timestamp := time.Now().Truncate(time.Second)
	r, err := d.insertPledge.Exec(itemID, userID, timestamp.Unix())
	if err != nil {
		return nil, err
	}
	lastID, err := r.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("no id returned for pledge for item %d for user %d: %v", itemID, userID, err)
	}
	return &types.Pledge{PledgeID: int(lastID), UserID: userID, Timestamp: timestamp}, nil
}

func (d *DB) ListUserPledges(userID string) ([]types.Pledge, error) {
	rows, err := d.Query("select pledgeID, userID, make, model, companyID, name, timestamp from userPledges where userID = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("ListUserPledges: unable to query: %v", err)
	}
	result := make([]types.Pledge, 0, 5)
	for rows.Next() {
		var p types.Pledge
		var timestamp int64
		err = rows.Scan(&p.PledgeID, &p.UserID, &p.Item.Make, &p.Item.Model, &p.Item.Company.ID, &p.Item.Company.Name, &timestamp)
		if err != nil {
			return result, fmt.Errorf("ListUserPledges: unable to parse result row: %v", err)
		}
		p.Timestamp = time.Unix(timestamp, 0)
		result = append(result, p)
	}
	return result, nil
}

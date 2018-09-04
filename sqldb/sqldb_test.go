package sqldb

import (
	"database/sql"
	"fmt"
	"psychic-rat/types"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestCreateDB(t *testing.T) {
	initDB(t)
}

func TestAddCompany(t *testing.T) {
	mock := mockDB{
		execs: []expectedExecStmt{
			{
				insert: &insertStmt{table: "companies",
					columns: map[string]interface{}{"name": "testco1"},
				},
				insertId: 1234},
		},
		t: t,
	}

	db := DB{mock}
	co, err := db.AddCompany(types.Company{Name: "testco1"})
	if err != nil {
		t.Fatal(err)
	}
	if 1234 != co.ID {
		t.Fatalf("expected ID 1234, got %v", co.ID)
	}
}

type expectedExecStmt struct {
	insert       *insertStmt
	insertId     int64
	rowsAffected int64
}

type insertStmt struct {
	table   string
	columns map[string]interface{}
}

func (m expectedExecStmt) LastInsertId() (int64, error) {
	return m.insertId, nil
}

func (m expectedExecStmt) RowsAffected() (int64, error) {
	return m.rowsAffected, nil
}

type mockDB struct {
	execsDone int
	execs     []expectedExecStmt
	t         *testing.T
}

func (m mockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	exec := m.execs[m.execsDone]
	m.execsDone++
	if exec.insert != nil {
		checkExecInsert(m.t, exec.insert, query, args)
		return exec, nil
	}
	return nil, fmt.Errorf("not able to match mock exec for query %v", query)
}

func checkExecInsert(t *testing.T, insert *insertStmt, query string, args []interface{}) {
	t.Helper()

	re := regexp.MustCompile("insert into (.*)\\((.*)\\) values\\((.*)\\)")
	results := re.FindStringSubmatch(query)
	if len(results) != 4 {
		t.Fatalf("could not match query: %v", query)
	}

	table := results[1]
	columns := strings.Split(results[2], ",")
	values := strings.Split(results[3], ",")

	if insert.table != table {
		t.Fatalf("wrong table, expected %v, got %v", insert.table, table)
	}
	if len(columns) != len(values) {
		t.Fatalf("wrong number of values, expected %v, got %v", len(columns), len(values))
	}
	if len(insert.columns) != len(columns) {
		t.Fatalf("column counts don't match, expected %v, got %v", insert.columns, columns)
	}
	for i, col := range columns {
		expv, exists := insert.columns[col]
		if !exists {
			t.Fatalf("unexpected column: %v", col)
		}
		if values[i] != "?" {
			t.Fatalf("got non-variable parameter, expected ?, got %v", values[i])
		}
		if !reflect.DeepEqual(expv, args[i]) {
			t.Fatalf("types & values don't match, expected %v of type %T, got %v of type %T", expv, expv, args[i], args[i])
		}

	}
}

func (m mockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	panic("not implemented")
}

func (m mockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	panic("not implemented")
}

func TestListCompanies(t *testing.T) {
	db := initDB(t)

	companies, err := db.ListCompanies()
	if err != nil {
		t.Fatal(err)
	}

	for i, c := range testCos {
		if c != companies[i].Name && i != companies[i].ID {
			t.Fatalf("company did not match, expected %v got %v", c, companies[i])
		}
	}

}

func TestGetCompanyById(t *testing.T) {
	db := initDB(t)

	id := 1
	c, err := db.GetCompany(id)
	if err != nil {
		t.Fatal(err)
	}
	if c.ID != id {
		t.Fatalf("wanted id %v, got %v", id, c.ID)
	}
}

func TestListItems(t *testing.T) {
	db := initDB(t)
	items, err := db.ListItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(testItems) != len(items) {
		t.Fatalf("expected %v items, got %v items [%v]", len(testItems), len(items), items)
	}
	for i := range items {
		if reflect.DeepEqual(testItems[i], items[i]) {
			t.Fatalf("expected %v, got %v", testItems[i], items[i])
		}
	}
}

func TestNewItems(t *testing.T) {
	db := initDB(t)
	newItems := initNewItems(db, t)
	ns, err := db.ListNewItems()
	if err != nil {
		t.Fatal(err)
	}
	// TODO: expected results template
	if len(newItems) != len(ns) {
		t.Fatalf("expected %v items, got %v items [%v]", len(newItems), len(ns), ns)
	}
	for i := range ns {
		if !reflect.DeepEqual(newItems[i], ns[i]) {
			t.Fatalf("expected item %v but got %v", newItems[i], ns[i])
		}
	}
}

func TestNewUser(t *testing.T) {
	db := initDB(t)

	for _, u := range testUsers {
		user, err := db.GetUser(u.ID)
		if err != nil {
			t.Fatal("error from db.GetUser:" + err.Error())
		}
		if reflect.DeepEqual(u, user) {
			t.Fatalf("expected %v got %v", u, user)
		}
	}
}

func TestGetItem(t *testing.T) {
	db := initDB(t)
	initCurrencies(db, t)
	ids := initItems(db, t)
	initCompanies(db, t)

	for i := range testItems {
		item, err := db.GetItem(ids[i])
		if err != nil {
			t.Fatal(err)
		}
		if reflect.DeepEqual(testItems[i], item) {
			t.Fatalf("expected %v, got %v", testItems[i], item)
		}
	}
}

func TestAddNewItem(t *testing.T) {
	db := initDB(t)
	initCurrencies(db, t)

	newItem := &types.NewItem{Make: "newthing", Model: "newmodel", Company: "newco", UserID: "test1", CurrencyID: 1, Value: 100}
	newItem, err := db.AddNewItem(*newItem)
	if err != nil {
		t.Fatal(err)
	}

	var n types.NewItem
	var timestamp int64
	err = db.QueryRow("select id, make, model, company, userID, currencyID, currencyValue, timestamp from newItems where id=?", newItem.ID).Scan(&n.ID, &n.Make, &n.Model, &n.Company, &n.UserID, &n.CurrencyID, &n.Value, &timestamp)
	if err != nil {
		t.Fatal(err)
	}
	n.Timestamp = time.Unix(timestamp, 0)
	if time.Since(n.Timestamp) > time.Second {
		t.Fatalf("timestamp problem, new is %v: ", n.Timestamp)
	}
	if !reflect.DeepEqual(*newItem, n) {
		t.Fatalf("expected %v, got back %v", newItem, n)
	}
}

func TestCurrencies(t *testing.T) {
	db := initDB(t)
	initCurrencies(db, t)

	currencies, err := db.ListCurrencies()
	if err != nil {
		t.Fatal(err)
	}

	for i, c := range testCurrencies {
		if c != currencies[i] {
			t.Fatalf("currency did not match, expected %v got %v", c, currencies[i])
		}
	}

}

func TestAddNewPledge(t *testing.T) {
	db := initDB(t)
	initCompanies(db, t)
	initItems(db, t)

	p, err := db.AddPledge(1, "user001", 100)
	if err != nil {
		t.Fatal(err)
	}

	if p.UserID != "user001" || p.USDValue != 100 || p.Item.ID != 1 {
		t.Fatalf("failed to add pledge, got %v", p)
	}
}

package sqldb

import (
	"psychic-rat/types"
	"reflect"
	"testing"
	"time"
)

func TestCreateDB(t *testing.T) {
	initDB(t)
}

func TestAddCompany(t *testing.T) {
	mock := NewMockDB(t)
	mock.ExecExpectation(NewExec("companies").
		WithColumnValue("name", "testco1").
		WithInsertId(1234))

	db := DB{mock}
	co, err := db.AddCompany(types.Company{Name: "testco1"})
	if err != nil {
		t.Fatal(err)
	}
	if 1234 != co.ID {
		t.Fatalf("expected ID 1234, got %v", co.ID)
	}
}

func TestListCompanies1(t *testing.T) {
	mock := NewMockDB(t)
	mock.QueryExpectation(NewQuery("companies").
		WithColumns("id", "name").
		WithResultsRow(1, "testco1").
		WithResultsRow(2, "testco2"))

	db := DB{mock}
	cos, err := db.ListCompanies()
	if err != nil {
		t.Fatal(err)
	}
	if len(cos) == 0 {
		t.Fatal("no companies returned!")
	}
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
	p, err := db.AddPledge(1, "user001", 100)
	if err != nil {
		t.Fatal(err)
	}

	if p.UserID != "user001" || p.USDValue != 100 {
		t.Fatalf("failed to add pledge, got %v", p)
	}
}

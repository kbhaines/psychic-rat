package sqldb

import (
	"psychic-rat/types"
	"reflect"
	"testing"
	"time"
)

func TestCreateDB(t *testing.T) {
	db := initDB(t)
	initCompanies(db, t)
}

func TestGetCompanyById(t *testing.T) {
	db := initDB(t)
	initCompanies(db, t)

	id := 1
	c, err := db.GetCompany(id)
	if err != nil {
		t.Fatal(err)
	}
	if c.ID != id {
		t.Fatalf("wanted id %v, got %v", id, c.ID)
	}
	if c.Name != testCos[0] {
		t.Fatalf("wanted name %v, got %v", testCos[0], c.Name)
	}
}

func TestListItems(t *testing.T) {
	db := initDB(t)
	initCompanies(db, t)
	initCurrencies(db, t)
	initItems(db, t)
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
	initCurrencies(db, t)
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
	initUsers(db, t)

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

package sqldb

import (
	"os"
	"psychic-rat/types"
	"reflect"
	"testing"
)

var (
	testCos = []string{"testco1", "testco2", "testco3"}

	testItems = []types.Item{}

	newItems = []types.NewItem{
		types.NewItem{Id: 1, UserID: 1, IsPledge: true, Make: "newPhone", Model: "newMod",
			Company: "co1", CompanyID: 1},
	}
)

func TestCreateDB(t *testing.T) {
	db := initDB(t)
	defer os.Remove("test.db")
	initCompanies(db, t)

	cs, err := db.GetCompanies()
	if err != nil {
		t.Fatal(err)
	}
	if len(testCos) != len(cs.Companies) {
		t.Fatalf("expected %v records but got %v\nRecords:", len(testCos), len(cs.Companies), cs.Companies)
	}
	for i, c := range testCos {
		if cs.Companies[i].Name != c {
			t.Fatal("did not get record back")
		}
	}
}

func initDB(t *testing.T) *DB {
	t.Helper()
	db, err := NewDB("test1.db")
	if err != nil {
		t.Fatalf("could not init DB: %v", err)
	}
	return db
}

func initCompanies(db *DB, t *testing.T) {
	for _, c := range testCos {
		err := db.NewCompany(types.Company{Name: c})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetCompanyById(t *testing.T) {
	db := initDB(t)
	defer os.Remove("test.db")
	initCompanies(db, t)

	id := 1
	c, err := db.GetCompany(id)
	if err != nil {
		t.Fatal(err)
	}
	if c.Id != id {
		t.Fatalf("wanted id %v, got %v", id, c.Id)
	}
	if c.Name != testCos[0] {
		t.Fatalf("wanted name %v, got %v", testCos[0], c.Name)
	}
}

func TestListItems(t *testing.T) {
	db := initDB(t)
	defer os.Remove("test.db")
	items, err := db.ListItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(testItems) != len(items.Items) {
		t.Fatalf("expected %v items, got %v items [%v]", len(testItems), len(items.Items), items.Items)
	}
}

func TestNewItems(t *testing.T) {
	db := initDB(t)
	defer os.Remove("test.db")
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

func initNewItems(db *DB, t *testing.T) []types.NewItem {
	res := []types.NewItem{}
	for _, c := range newItems {
		n, err := db.AddNewItem(c)
		if err != nil {
			t.Fatal(err)
		}
		res = append(res, *n)
	}
	return res
}

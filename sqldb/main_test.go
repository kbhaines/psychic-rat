package sqldb

import (
	"psychic-rat/mdl"
	"psychic-rat/types"
	"testing"
)

var (
	testCos = []string{"testco1", "testco2", "testco3"}

	testItems = []types.Item{}
)

func TestCreateDB(t *testing.T) {
	db := initDB(t)
	cs, err := db.GetCompanies()
	if err != nil {
		t.Fatal(err)
	}
	if len(cs.Companies) != len(testCos) {
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
	db, err := NewDB("test.db")
	if err != nil {
		t.Fatalf("could not init DB: %v", err)
	}

	for _, c := range testCos {
		err = db.NewCompany(mdl.Company{Name: c})
		if err != nil {
			t.Fatal(err)
		}
	}
	return db
}

func TestGetCompanyById(t *testing.T) {
	db := initDB(t)
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
	items, err := db.ListItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(items.Items) != len(testItems) {
		t.Fatalf("expected %v items, got %v items [%v]", len(testItems), len(items.Items), items.Items)
	}
}

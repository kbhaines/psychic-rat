package sqldb

import (
	"psychic-rat/mdl"
	"testing"
)

func TestCreateDB(t *testing.T) {
	db, err := NewDB("test.db")
	if err != nil {
		t.Fatalf("could not init DB: %v", err)
	}

	testCos := []string{"testco1", "testco2", "testco3"}
	for _, c := range testCos {
		err = db.NewCompany(mdl.Company{Name: c})
		if err != nil {
			t.Fatal(err)
		}
	}
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

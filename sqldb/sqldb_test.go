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

var (
	cos = []types.Company{
		types.Company{ID: 1, Name: "testco1"},
		types.Company{ID: 2, Name: "testco2"},
		types.Company{ID: 3, Name: "testco3"},
		types.Company{ID: 4, Name: "testco4"},
	}

	items = []types.Item{
		types.Item{ID: 1, Make: "phone", Model: "xyz", Company: cos[0], USDValue: 100},
		types.Item{ID: 2, Make: "phone", Model: "133", Company: cos[0], USDValue: 100},
		types.Item{ID: 3, Make: "tablet", Model: "ab1", Company: cos[1], USDValue: 100},
		types.Item{ID: 4, Make: "tablet", Model: "xy1", Company: cos[1], USDValue: 100},
	}

	currencies = []types.Currency{
		types.Currency{ID: 1, Ident: "USD", ConversionToUSD: 1.0},
		types.Currency{ID: 2, Ident: "GBP", ConversionToUSD: 1.2},
		types.Currency{ID: 3, Ident: "EUR", ConversionToUSD: 1.4},
	}
)

func TestAddCompany(t *testing.T) {
	mock := NewMockDB(t).
		ExecExpectation(NewExec("companies").
			WithColumnValue("name", "testco1").
			WithInsertId(1234))

	db := DB{mock}
	co, err := db.AddCompany(types.Company{Name: "testco1"})
	mock.CheckAllExpectationsMet()
	if err != nil {
		t.Fatal(err)
	}
	if 1234 != co.ID {
		t.Fatalf("expected ID 1234, got %v", co.ID)
	}
}

func TestListCompanies(t *testing.T) {
	qe := NewQuery("companies").WithColumns("id", "name")
	for _, r := range cos {
		qe.WithResultsRow(r.ID, r.Name)
	}
	mock := NewMockDB(t).QueryExpectation(qe)

	db := DB{mock}
	listing, err := db.ListCompanies()
	if err != nil {
		t.Fatal(err)
	}
	mock.CheckAllExpectationsMet()

	if len(cos) != len(listing) {
		t.Errorf("expected %d companies, got %d", len(cos), len(listing))
	}
}

func TestGetCompanyById(t *testing.T) {
	co := cos[0]
	mock := NewMockDB(t).
		QueryExpectation(NewQuery("companies").
			WithColumns("id", "name").
			WithResultsRow(co.ID, co.Name))

	db := DB{mock}
	c, err := db.GetCompany(co.ID)
	mock.CheckAllExpectationsMet()
	if err != nil {
		t.Fatal(err)
	}
	if co != c {
		t.Fatalf("wanted company record %v, got %v", co, c)
	}
}

func TestGetCompanyByWrongId(t *testing.T) {
	co := cos[0]
	mock := NewMockDB(t).
		QueryExpectation(NewQuery("companies").
			WithColumns("id", "name").
			WithError("not found"))

	db := DB{mock}
	_, err := db.GetCompany(co.ID)
	mock.CheckAllExpectationsMet()
	if err == nil || err.Error() != "not found" {
		t.Fatalf("did not get expected error, got %v", err)
	}
}

func TestListItems(t *testing.T) {
	qe := NewQuery("itemsCompany").
		WithColumns("id", "make", "model", "companyID", "companyName", "usdValue")
	for _, i := range items {
		qe.WithResultsRow(i.ID, i.Make, i.Model, i.Company.ID, i.Company.Name, i.USDValue)
	}

	mock := NewMockDB(t).QueryExpectation(qe)
	db := DB{mock}
	its, err := db.ListItems()
	if err != nil {
		t.Fatal(err)
	}
	mock.CheckAllExpectationsMet()
	if len(items) != len(its) {
		t.Fatalf("expected %v items, got %v items [%v]", len(items), len(its), its)
	}
	for i := range its {
		if !reflect.DeepEqual(items[i], its[i]) {
			t.Fatalf("expected %v, got %v", items[i], its[i])
		}
	}
}

func TestGetItem(t *testing.T) {
	mock := NewMockDB(t)
	for _, i := range items {
		mock.QueryExpectation(NewQuery("itemsCompany").
			WithColumns("id", "make", "model", "companyID", "companyName", "usdValue").
			WithResultsRow(i.ID, i.Make, i.Model, i.Company.ID, i.Company.Name, i.USDValue))
	}

	db := DB{mock}

	for _, item := range items {
		it, err := db.GetItem(item.ID)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(item, it) {
			t.Fatalf("expected %v, got %v", item, it)
		}
	}
	mock.CheckAllExpectationsMet()
}

func TestListCurrencies(t *testing.T) {
	qe := NewQuery("currencies").
		WithColumns("id", "ident", "usdConversion")
	for _, c := range currencies {
		qe.WithResultsRow(c.ID, c.Ident, c.ConversionToUSD)
	}

	mock := NewMockDB(t).QueryExpectation(qe)
	db := DB{mock}

	curs, err := db.ListCurrencies()
	if err != nil {
		t.Fatal(err)
	}
	mock.CheckAllExpectationsMet()

	for i := range curs {
		if currencies[i] != curs[i] {
			t.Fatalf("currency did not match, expected %v got %v", currencies[i], curs[i])
		}
	}
}
func TestAddNewItems(t *testing.T) {
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

func TestListNewItems(t *testing.T) {
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

func TestAddUser(t *testing.T) {
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

func TestAddPledge(t *testing.T) {
	db := initDB(t)
	p, err := db.AddPledge(1, "user001", 100)
	if err != nil {
		t.Fatal(err)
	}

	if p.UserID != "user001" || p.USDValue != 100 {
		t.Fatalf("failed to add pledge, got %v", p)
	}
}

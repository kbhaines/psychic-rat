package sqldb

import (
	"psychic-rat/types"
	"reflect"
	"testing"
	"time"
)

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

	newItems = []types.NewItem{
		types.NewItem{ID: 1, UserID: "test1", IsPledge: true, Make: "newPhone", Model: "newMod", Company: "co1", CompanyID: 1, CurrencyID: 1, Value: 100, Timestamp: time.Unix(0, 0)},
		types.NewItem{ID: 2, UserID: "test2", IsPledge: true, Make: "newPhone", Model: "newMod", Company: "co1", CompanyID: 1, CurrencyID: 1, Value: 100, Timestamp: time.Unix(0, 0)},
		types.NewItem{ID: 3, UserID: "test3", IsPledge: true, Make: "newPhone", Model: "newMod", Company: "co1", CompanyID: 1, CurrencyID: 1, Value: 100, Timestamp: time.Unix(0, 0)},
	}

	users = []types.User{
		types.User{ID: "test1", FirstName: "user1", Email: "user1@user.com"},
		types.User{ID: "test2", FirstName: "user2", Email: "user2@user.com"},
		types.User{ID: "test3", FirstName: "user3", Email: "user3@user.com"},
		types.User{ID: "test4", FirstName: "user4", Email: "user4@user.com"},
		types.User{ID: "test5", FirstName: "user5", Email: "user5@user.com"},
		types.User{ID: "test6", FirstName: "user6", Email: "user6@user.com"},
		types.User{ID: "test7", FirstName: "user7", Email: "user7@user.com"},
		types.User{ID: "test8", FirstName: "user8", Email: "user8@user.com"},
		types.User{ID: "test9", FirstName: "user9", Email: "user9@user.com"},
		types.User{ID: "test10", FirstName: "user10", Email: "user10@user.com"},
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
func TestAddNewItem(t *testing.T) {
	ni := newItems[0]
	mock := NewMockDB(t)

	mock.QueryExpectation(NewQuery("currencies").
		WithColumns("id", "ident", "usdConversion").
		WithResultsRow(ni.CurrencyID, "USD", 2.0))

	mock.ExecExpectation(NewExec("newItems").
		WithColumnValue("userID", ni.UserID).
		WithColumnValue("isPledge", ni.IsPledge).
		WithColumnValue("make", ni.Make).
		WithColumnValue("model", ni.Model).
		WithColumnValue("company", ni.Company).
		WithColumnValue("companyID", ni.CompanyID).
		WithColumnValue("currencyID", ni.CurrencyID).
		WithColumnValue("currencyValue", ni.Value).
		WithColumnValue("timestamp", time.Now().Truncate(time.Second).Unix()).
		WithColumnValue("used", 0))

	db := DB{mock}
	_, err := db.AddNewItem(ni)
	if err != nil {
		t.Fatal(err)
	}
	mock.CheckAllExpectationsMet()

}

func TestListNewItems(t *testing.T) {
	qe := NewQuery("newItems").
		WithColumns("id", "userID", "isPledge", "make", "model", "company", "companyID", "currencyID", "currencyValue", "timestamp")

	for _, ni := range newItems {
		qe.WithResultsRow(ni.ID, ni.UserID, ni.IsPledge, ni.Make, ni.Model, ni.Company, ni.CompanyID, ni.CurrencyID, ni.Value, int64(0))
	}

	mock := NewMockDB(t).QueryExpectation(qe)

	db := DB{mock}
	ns, err := db.ListNewItems()
	if err != nil {
		t.Fatal(err)
	}
	mock.CheckAllExpectationsMet()
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
	user := users[0]
	mock := NewMockDB(t).ExecExpectation(NewExec("users").
		WithColumnValue("id", user.ID).
		WithColumnValue("fullName", user.Fullname).
		WithColumnValue("firstName", user.FirstName).
		WithColumnValue("country", user.Country).
		WithColumnValue("email", user.Email).
		WithColumnValue("isAdmin", user.IsAdmin).
		WithInsertId(1234))

	db := DB{mock}
	err := db.AddUser(user)
	if err != nil {
		t.Fatal(err)
	}
	mock.CheckAllExpectationsMet()
}

func TestGetUser(t *testing.T) {
	user := users[0]
	mock := NewMockDB(t).QueryExpectation(NewQuery("users").
		WithColumns("id", "fullname", "firstName", "country", "email", "isAdmin").
		WithResultsRow(user.ID, user.Fullname, user.FirstName, user.Country, user.Email, user.IsAdmin))

	db := DB{mock}
	u, err := db.GetUser(user.ID)
	if err != nil {
		t.Fatal("error from db.GetUser:" + err.Error())
	}
	mock.CheckAllExpectationsMet()
	if reflect.DeepEqual(u, user) {
		t.Fatalf("expected %v got %v", u, user)
	}
}

func TestAddPledge(t *testing.T) {
	mock := NewMockDB(t).ExecExpectation(NewExec("pledges").
		WithInsertId(1234).
		WithColumnValue("itemID", 1).
		WithColumnValue("userID", "user001").
		WithColumnValue("usdValue", 100).
		WithColumnValue("timestamp", time.Now().Truncate(time.Second).Unix()))

	db := DB{mock}
	p, err := db.AddPledge(1, "user001", 100)
	if err != nil {
		t.Fatal(err)
	}
	mock.CheckAllExpectationsMet()

	if p.UserID != "user001" || p.USDValue != 100 {
		t.Fatalf("failed to add pledge, got %v", p)
	}
}

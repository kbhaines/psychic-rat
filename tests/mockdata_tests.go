package tests

import (
	"psychic-rat/sqldb"
	"psychic-rat/types"
	"testing"
)

type DB = sqldb.DB

var (
	testUsers = []types.User{
		types.User{ID: "test1", Fullname: "user1 full", Email: "user1@user.com"},
		types.User{ID: "test2", Fullname: "user2 full", Email: "user2@user.com"},
		types.User{ID: "test3", Fullname: "user3 full", Email: "user3@user.com"},
		types.User{ID: "test4", Fullname: "user4 full", Email: "user4@user.com"},
		types.User{ID: "test5", Fullname: "user5 full", Email: "user5@user.com"},
		types.User{ID: "test6", Fullname: "user6 full", Email: "user6@user.com"},
		types.User{ID: "test7", Fullname: "user7 full", Email: "user7@user.com"},
		types.User{ID: "test8", Fullname: "user8 full", Email: "user8@user.com"},
		types.User{ID: "test9", Fullname: "user9 full", Email: "user9@user.com"},
		types.User{ID: "test10", Fullname: "user10 full", Email: "user10@user.com"},
		types.User{ID: "admin", Fullname: "Admin", Email: "admin@admin.com", IsAdmin: true},
	}

	testCompanies = []types.Company{
		types.Company{1, "testco1"},
		types.Company{2, "testco2"},
		types.Company{3, "testco3"},
	}

	testItems = []types.Item{
		types.Item{ID: 0, Make: "phone", Model: "xyz", Company: testCompanies[0], USDValue: 100},
		types.Item{ID: 0, Make: "phone", Model: "133", Company: testCompanies[0], USDValue: 100},
		types.Item{ID: 0, Make: "tablet", Model: "ab1", Company: testCompanies[1], USDValue: 100},
		types.Item{ID: 0, Make: "tablet", Model: "xy1", Company: testCompanies[1], USDValue: 100},
	}

	testNewItems = []types.NewItem{
		types.NewItem{ID: 1, UserID: "test1", IsPledge: true, Make: "newPhone1", Model: "newMod1", Company: "co1", CompanyID: 1, CurrencyID: 1, Value: 100},
		types.NewItem{ID: 2, UserID: "test2", IsPledge: true, Make: "newPhone2", Model: "newMod2", Company: "co2", CompanyID: 1, CurrencyID: 1, Value: 100},
		types.NewItem{ID: 3, UserID: "test3", IsPledge: true, Make: "newPhone3", Model: "newMod3", Company: "co3", CompanyID: 1, CurrencyID: 1, Value: 100},
		types.NewItem{ID: 4, UserID: "test4", IsPledge: true, Make: "newPhone4", Model: "newMod4", Company: "co4", CompanyID: 1, CurrencyID: 1, Value: 100},
		types.NewItem{ID: 5, UserID: "test5", IsPledge: true, Make: "newPhone5", Model: "newMod5", Company: "co5", CompanyID: 1, CurrencyID: 1, Value: 100},
		types.NewItem{ID: 6, UserID: "test6", IsPledge: true, Make: "newPhone6", Model: "newMod6", Company: "co6", CompanyID: 1, CurrencyID: 1, Value: 100},
	}

	testCurrencies = []types.Currency{
		types.Currency{ID: 1, Ident: "USD", ConversionToUSD: 1.0},
		types.Currency{ID: 2, Ident: "GBP", ConversionToUSD: 2},
		types.Currency{ID: 3, Ident: "EUR", ConversionToUSD: 4},
	}
)

func initDB(t *testing.T) *sqldb.DB {
	t.Helper()
	sql, err := sqldb.NewSqliteDB(":memory:")
	db := sqldb.NewDB(sql)
	if err != nil {
		t.Fatal(err)
	}
	initCompanies(db, t)
	initUsers(db, t)
	initCurrencies(db, t)
	//initNewItems(db,t)
	initItems(db, t)
	return db
}

func initCompanies(db *DB, t *testing.T) {
	for _, c := range testCompanies {
		_, err := db.AddCompany(c)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func initItems(db *DB, t *testing.T) []int {
	t.Helper()
	ids := []int{}
	for _, item := range testItems {
		i, err := db.AddItem(item)
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, i.ID)
	}
	return ids
}
func initNewItems(db *DB, t *testing.T) []types.NewItem {
	res := []types.NewItem{}
	for _, c := range testNewItems {
		n, err := db.AddNewItem(c)
		if err != nil {
			t.Fatal(err)
		}
		res = append(res, *n)
	}
	return res
}

func initUsers(db *DB, t *testing.T) {
	t.Helper()
	for _, u := range testUsers {
		err := db.AddUser(u)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func initCurrencies(db *DB, t *testing.T) {
	t.Helper()
	for _, c := range testCurrencies {
		_, err := db.AddCurrency(c)
		if err != nil {
			t.Fatal(err)
		}
	}
}

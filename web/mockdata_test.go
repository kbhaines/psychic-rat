package web

import (
	"psychic-rat/types"
	"psychic-rat/sqldb"
	"testing"
)

type DB = sqldb.DB

var (
	testUsers = []types.User{
		types.User{Id: "test1", Fullname: "user1 full", Email: "user1@user.com"},
		types.User{Id: "test2", Fullname: "user2 full", Email: "user2@user.com"},
		types.User{Id: "test3", Fullname: "user3 full", Email: "user3@user.com"},
		types.User{Id: "test4", Fullname: "user4 full", Email: "user4@user.com"},
		types.User{Id: "test5", Fullname: "user5 full", Email: "user5@user.com"},
		types.User{Id: "test6", Fullname: "user6 full", Email: "user6@user.com"},
		types.User{Id: "test7", Fullname: "user7 full", Email: "user7@user.com"},
		types.User{Id: "test8", Fullname: "user8 full", Email: "user8@user.com"},
		types.User{Id: "test9", Fullname: "user9 full", Email: "user9@user.com"},
		types.User{Id: "test10", Fullname: "user10 full", Email: "user10@user.com"},
		types.User{Id: "admin", Fullname: "Admin", Email: "admin@admin.com", IsAdmin: true},
	}

	testCos = []string{"testco1", "testco2", "testco3"}

	companies = []types.Company{
		types.Company{1, "testco1"},
		types.Company{2, "testco2"},
	}

	testItems = []types.Item{
		types.Item{Id: 0, Make: "phone", Model: "xyz", Company: companies[0]},
		types.Item{Id: 0, Make: "phone", Model: "133", Company: companies[0]},
		types.Item{Id: 0, Make: "tablet", Model: "ab1", Company: companies[1]},
		types.Item{Id: 0, Make: "tablet", Model: "xy1", Company: companies[1]},
	}

	newItems = []types.NewItem{
		types.NewItem{Id: 1, UserID: "test1", IsPledge: true, Make: "newPhone", Model: "newMod", Company: "co1", CompanyID: 1},
		types.NewItem{Id: 2, UserID: "test2", IsPledge: true, Make: "newPhone", Model: "newMod", Company: "co2", CompanyID: 1},
		types.NewItem{Id: 3, UserID: "test3", IsPledge: true, Make: "newPhone", Model: "newMod", Company: "co3", CompanyID: 1},
		types.NewItem{Id: 4, UserID: "test4", IsPledge: true, Make: "newPhone", Model: "newMod", Company: "co4", CompanyID: 1},
		types.NewItem{Id: 5, UserID: "test5", IsPledge: true, Make: "newPhone", Model: "newMod", Company: "co5", CompanyID: 1},
		types.NewItem{Id: 6, UserID: "test6", IsPledge: true, Make: "newPhone", Model: "newMod", Company: "co6", CompanyID: 1},
	}
)

func initDB(t *testing.T) {
	t.Helper()
	db,err := sqldb.NewDB("test.dat")
	if err != nil {
		t.Fatal(err)
	}
	apis = API{
		Company: db,
		Item:    db,
		NewItem: db,
		Pledge:  db,
		User:    db,
	}
	initCompanies(db,t)
	initUsers(db,t)
	//initNewItems(db,t)
	initItems(db,t)
}

func initCompanies(db *DB, t *testing.T) {
	for _, c := range testCos {
		_, err := db.NewCompany(types.Company{Name: c})
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
		ids = append(ids, i.Id)
	}
	return ids
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

func initUsers(db *DB, t *testing.T) {
	t.Helper()
	for _, u := range testUsers {
		err := db.CreateUser(u)
		if err != nil {
			t.Fatal(err)
		}
	}
}

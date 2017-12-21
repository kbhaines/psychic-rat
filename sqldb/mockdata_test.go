package sqldb

import (
	"psychic-rat/types"
	"testing"
)

var (
	testUsers = []types.User{
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

	testCos = []string{"testco1", "testco2", "testco3"}

	companies = []types.Company{
		types.Company{1, "testco1"},
		types.Company{2, "testco2"},
	}

	testItems = []types.Item{
		types.Item{ID: 0, Make: "phone", Model: "xyz", Company: companies[0]},
		types.Item{ID: 0, Make: "phone", Model: "133", Company: companies[0]},
		types.Item{ID: 0, Make: "tablet", Model: "ab1", Company: companies[1]},
		types.Item{ID: 0, Make: "tablet", Model: "xy1", Company: companies[1]},
	}

	newItems = []types.NewItem{
		types.NewItem{UserID: "test1", IsPledge: true, Make: "newPhone", Model: "newMod", Company: "co1", CompanyID: 1},
		types.NewItem{UserID: "test2", IsPledge: true, Make: "newPhone", Model: "newMod", Company: "co1", CompanyID: 1},
		types.NewItem{UserID: "test3", IsPledge: true, Make: "newPhone", Model: "newMod", Company: "co1", CompanyID: 1},
	}
)

func initDB(t *testing.T) *DB {
	t.Helper()
	db, err := NewDB("test.db")
	if err != nil {
		t.Fatalf("could not init DB: %v", err)
	}
	return db
}

func initCompanies(db *DB, t *testing.T) {
	for _, c := range testCos {
		_, err := db.AddCompany(types.Company{Name: c})
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
		err := db.AddUser(u)
		if err != nil {
			t.Fatal(err)
		}
	}
}

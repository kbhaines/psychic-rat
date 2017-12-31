package datagen

import (
	"fmt"
	"log"
	"psychic-rat/sqldb"
	"psychic-rat/types"
	"time"
)

func Generate(db *sqldb.DB, totalSize int) error {
	defer timeTrack(time.Now(), "Generate DB")

	numCompanies := totalSize / 1000
	numUsers := totalSize / 10
	numItem := totalSize / 100

	genCos := func() error {
		defer timeTrack(time.Now(), "Gen cos")
		for c := 0; c < numCompanies; c++ {
			if _, err := db.AddCompany(generateCompany(c)); err != nil {
				return err
			}
		}
		return nil
	}

	genUsers := func() error {
		defer timeTrack(time.Now(), "Gen users")
		for u := 0; u < numUsers; u++ {
			if err := db.AddUser(generateUser(u)); err != nil {
				return err
			}
		}
		return nil
	}
	genItems := func() error {
		defer timeTrack(time.Now(), "Gen items")
		for i := 0; i < numItem; i++ {
			if _, err := db.AddItem(generateItem(i, numCompanies)); err != nil {
				return err
			}
		}
		return nil
	}
	genNewItems := func() error {
		defer timeTrack(time.Now(), "Gen new items")
		for i := 0; i < numItem; i++ {
			if _, err := db.AddNewItem(generateNewItem(i, numItem)); err != nil {
				return err
			}
		}
		return nil
	}

	runOrPanic(genCos)
	runOrPanic(genUsers)
	runOrPanic(genItems)
	runOrPanic(genNewItems)
	return nil
}

func runOrPanic(f func() error) {
	if err := f(); err != nil {
		panic(fmt.Sprintf("failed %v: %v", f, err))
	}
}

var spf = fmt.Sprintf

func generateCompany(c int) types.Company {
	return types.Company{Name: spf("company%03d", c)}
}

func generateItem(i, maxCompanyId int) types.Item {
	company := types.Company{ID: (i % maxCompanyId) + 1}
	return types.Item{Make: spf("make%03d", i), Model: spf("model%03d", i), Company: company}
}

func generateNewItem(i, maxCompanyId int) types.NewItem {
	return types.NewItem{UserID: spf("user%03d", i), Make: spf("newmake%03d", i), Model: spf("newmodel%03d", i), Company: spf("newco%03d", i), IsPledge: true}
}

func generateUser(u int) types.User {
	return types.User{ID: spf("user%03d", u), FirstName: spf("User%03d", u), Fullname: spf("User%03d Fullname", u), Email: spf("user%03d@domain.com", u), IsAdmin: u == 42}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

package datagen

import (
	"fmt"
	"log"
	"psychic-rat/mdl"
	"psychic-rat/sqldb"
	"psychic-rat/types"
	"time"
)

func Generate(db *sqldb.DB, totalSize int) error {
	defer timeTrack(time.Now(), "Generate DB")

	numCompanies := totalSize / 1000
	numUsers := totalSize / 10
	numItem := totalSize / 100

	genCos := func() {
		defer timeTrack(time.Now(), "Gen cos")
		for c := 0; c < numCompanies; c++ {
			err := db.NewCompany(generateCompany(c))
			if err != nil {
				return
			}
		}
	}
	genUsers := func() {
		defer timeTrack(time.Now(), "Gen users")
		for u := 0; u < numUsers; u++ {
			err := db.CreateUser(generateUser(u))
			if err != nil {
				return
			}
		}
	}
	genItems := func() {
		defer timeTrack(time.Now(), "Gen items")
		for i := 0; i < numItem; i++ {
			_, err := db.AddItem(generateItem(i, numCompanies))
			if err != nil {
				return
			}
		}
	}

	genCos()
	genUsers()
	genItems()
	return nil
}

var spf = fmt.Sprintf

func generateCompany(c int) types.Company {
	return types.Company{Name: spf("company%03d", c)}
}

func generateItem(i, maxCompanyId int) types.Item {
	company := types.Company{Id: i % maxCompanyId}
	return types.Item{Make: spf("make%03d", i), Model: spf("model%03d", i), Company: company}
}

func generateUser(u int) mdl.User {
	return mdl.User{Id: spf("user%03d", u), FirstName: spf("User%03d", u), Fullname: spf("User%03d Fullname", u), Email: spf("user%03d@domain.com", u)}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

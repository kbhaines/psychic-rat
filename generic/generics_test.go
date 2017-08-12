package generic

import "testing"

type person struct {
	Name string
	Age  int
}

type car struct {
	Make  string
	Model string
}

var people = [...]person{
	person{"Kevin", 47},
	person{"Marcelle", 21},
	person{"Loen", 39},
}

var cars = [...]car{
	car{"Rover", "620si"},
	car{"Mercedes", "220cdi"},
}

func createPerson(db GenericDb, t *testing.T, ids ...int) []GenericId {
	results := make([]GenericId, len(ids))
	j := 0
	for i := range ids {
		id, err := db.Create(people[i])
		if err != nil {
			t.Errorf("could not create %v, got %v", people[0], err)
		}
		results[j] = id
		j++
	}
	return results
}

func TestCreate(t *testing.T) {
	db := MakeDb()
	id := createPerson(db, t, 0)[0]
	personOut := checkGetPerson(db, id, t)
	checkExpectedPerson(t, people[0], personOut)
}

func TestMultiCreate(t *testing.T) {
	db := MakeDb()
	ids := createPerson(db, t, 0, 1, 2)
	j := 0
	for i := range ids {
		personOut := checkGetPerson(db, i, t)
		checkExpectedPerson(t, people[j], personOut)
		j++
	}
}

func checkExpectedPerson(t *testing.T, expected person, actual person) {
	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	t.Logf("%s: got %v, yay!", t.Name(), actual)
}

func checkGetPerson(db GenericDb, i interface{}, t *testing.T) person {
	p, err := db.Get(i)
	if err != nil {
		t.Errorf("did not retrieve person")
	}
	return p.(person)
}

func TestNotFound(t *testing.T) {
	db := MakeDb()
	createPerson(db, t, 0)
	_, err := db.Get(1)
	if err == nil {
		t.Errorf("error not found")
	}
}

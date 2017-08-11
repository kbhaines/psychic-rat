package entities

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
	p, err := db.Get(id)
	if err != nil {
		t.Errorf("did not retrieve person")
	}
	personOut := p.(person)
	if personOut != people[0] {
		t.Errorf("expected %v, got back %v", people[0], personOut)
	}
	t.Logf("got %v, yay!", personOut)
}

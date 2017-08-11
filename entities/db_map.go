package entities

import "errors"

type ItemDb struct {
	items map[ItemId]Item
}


func (db ItemDb) Create(item Item) (ItemId, error) {
	newId := ItemId(len(db.items))
	db.items[newId] = item
	return newId, nil
}

func (db ItemDb) GetById(id ItemId) (Item, error) {
	item, exists := db.items[id]
	if ! exists {
		return Item{}, errors.New("not found")
	}
	return item, nil
}

func (db ItemDb) List() []ItemId {
	items := make([]ItemId, len(db.items))
	j := 0
	for i := range db.items {
		items[j] = i
		j++
	}
	return items
}

type PublicUserDb struct {
	users map[UserId]PublicUser
}

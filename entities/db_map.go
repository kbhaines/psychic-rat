package entities

type ItemDb struct {
	items map[ItemId]Item
}

type GenericId interface{}
type Generic interface{}

type GenericMap struct {
	items map[GenericId]Generic
}

func (d GenericMap) Create(item Generic) (GenericId, error) {
	newId := len(d.items)
	d.items[newId] = item
	return newId, nil
}

func (d GenericMap) Get(id GenericId) (Generic, error) {
	item, exists := d.items[id]
	if ! exists {
		return nil, error("not found")
	}
	return item, nil
}

func (d GenericMap) ListIds() ([]GenericId) {
	ids := make([]GenericId, len(d.items))
	j := 0
	for i := range d.items {
		ids[j] = i
		j++
	}
	return ids
}

func (db ItemDb) Create(item Item) (ItemId, error) {
	newId := ItemId(len(db.items))
	db.items[newId] = item
	return newId, nil
}

func (db ItemDb) GetById(id ItemId) (Item, error) {
	item, exists := db.items[id]
	if ! exists {
		return Item{}, error("not found")
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

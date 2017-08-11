package entities

import (
	"errors"
	"reflect"
)

type GenericId interface{}
type Generic struct {
	T reflect.Type
	V interface{}
}

type GenericDb struct {
	items map[GenericId]interface{}
}

func MakeDb() GenericDb {
	return GenericDb{make(map[GenericId]interface{})}
}

func (d GenericDb) Create(item interface{}) (GenericId, error) {
	newId := len(d.items)
	d.items[newId] = item
	return newId, nil
}

func (d GenericDb) Get(id GenericId) (interface{}, error) {
	item, exists := d.items[id]
	if ! exists {
		return nil, errors.New("not found")
	}
	return item, nil
}

func (d GenericDb) ListIds() ([]GenericId) {
	ids := make([]GenericId, len(d.items))
	j := 0
	for i := range d.items {
		ids[j] = i
		j++
	}
	return ids
}

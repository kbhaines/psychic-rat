package itemrepo

import (
	"psychic-rat/mdl/item"
	"psychic-rat/mdl/company"
	"psychic-rat/repo"
	"errors"
	"fmt"
)

func GetItemRepoMapImpl() repo.Items {
	return itemRepo
}

var itemRepo repo.Items = &repoMap{make(map[item.Id]item.Record)}

type repoMap struct {
	records map[item.Id]item.Record
}

type record struct {
	id item.Id
	item.Record
}

func (r record) RepoId() item.Id {
	return r.id
}

func (r *repoMap) Create(i item.Record) error {
	if _, found := r.records[i.Id()]; found {
		return fmt.Errorf("item %v already exists", i.Id())
	}
	r.records[i.Id()] = i
	return nil
}

func (r *repoMap) GetById(id item.Id) (item.Record, error) {
	item, found := r.records[id]
	if !found {
		return nil, errors.New("not found")
	}
	return record{id, item}, nil
}

func (r *repoMap) Update(id item.Id, item item.Record) {
	panic("implement me")
}

func (r *repoMap) List() []item.Id {
	itemIds := make([]item.Id, len(r.records))
	i := 0
	for id := range r.records {
		itemIds[i] = id
		i++
	}
	return itemIds
}

func (r *repoMap) GetAllByCompany(companyId company.Id) (items []item.Record) {
	for id, r := range r.records {
		if r.Company() == companyId {
			items = append(items, record{id, r})
		}
	}
	return items
}

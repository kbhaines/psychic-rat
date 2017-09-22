package repo

import (
	"errors"
	"fmt"
	"psychic-rat/mdl"
)

func getItemRepoMapImpl() Items {
	return itemRepo
}

var itemRepo Items = &itemRepoMap{make(map[mdl.Id]mdl.ItemRecord)}

type itemRepoMap struct {
	records map[mdl.Id]mdl.ItemRecord
}

type record struct {
	id mdl.Id
	mdl.ItemRecord
}

func (r record) RepoId() mdl.Id {
	return r.id
}

func (r *itemRepoMap) Create(i mdl.ItemRecord) error {
	if _, found := r.records[i.Id]; found {
		return fmt.Errorf("item %v already exists", i.Id)
	}
	r.records[i.Id] = i
	return nil
}

func (r *itemRepoMap) GetById(id mdl.Id) (*mdl.ItemRecord, error) {
	item, found := r.records[id]
	if !found {
		return nil, errors.New("not found")
	}
	return &item, nil
}

func (r *itemRepoMap) Update(id mdl.Id, item mdl.ItemRecord) {
	panic("implement me")
}

func (r *itemRepoMap) List() []mdl.ItemRecord {
	items := make([]mdl.ItemRecord, len(r.records))
	i := 0
	for _, item := range r.records {
		items[i] = item
		i++
	}
	return items
}

func (r *itemRepoMap) GetAllByCompany(companyId mdl.Id) (items []mdl.ItemRecord) {
	for _, r := range r.records {
		if r.CompanyId == companyId {
			items = append(items, r)
		}
	}
	return items
}

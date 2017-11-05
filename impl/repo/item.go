package repo

import (
	"errors"
	"fmt"
	"psychic-rat/mdl"
)

var itemRepo = &itemRepoMap{make(map[mdl.ID]mdl.Item)}

type itemRepoMap struct {
	records map[mdl.ID]mdl.Item
}

type record struct {
	id mdl.ID
	mdl.Item
}

func (r record) RepoId() mdl.ID {
	return r.id
}

func (r *itemRepoMap) Create(i mdl.Item) error {
	if _, found := r.records[i.Id]; found {
		return fmt.Errorf("item %v already exists", i.Id)
	}
	r.records[i.Id] = i
	return nil
}

func (r *itemRepoMap) GetById(id mdl.ID) (*mdl.Item, error) {
	item, found := r.records[id]
	if !found {
		return nil, errors.New("not found")
	}
	return &item, nil
}

func (r *itemRepoMap) Update(id mdl.ID, item mdl.Item) {
	panic("implement me")
}

func (r *itemRepoMap) List() []mdl.Item {
	items := make([]mdl.Item, len(r.records))
	i := 0
	for _, item := range r.records {
		items[i] = item
		i++
	}
	return items
}

func (r *itemRepoMap) GetAllByCompany(companyId mdl.ID) (items []mdl.Item) {
	for _, r := range r.records {
		if r.CompanyID == companyId {
			items = append(items, r)
		}
	}
	return items
}

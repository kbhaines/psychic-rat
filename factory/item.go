package factory

import "errors"
import (
	"psychic-rat/mdl/item"
	"psychic-rat/mdl/company"
)

func GetItemRepo() item.Repo {
	return itemRepo
}

// declare that we implement Repo interface
var itemRepo item.Repo = &repoMap{ make(map[item.Id]item.Record) }

type repoMap struct {
	records map[item.Id]item.Record
}

func (repo *repoMap) Create(i item.Record) (item.Id, error) {
	repo.records[i.Id()] = i
	return i.Id(), nil
}

func (repo *repoMap) GetById(id item.Id) (item.Record, error) {
	item, found := repo.records[id]
	if !found {
		return nil, errors.New("not found")
	}
	return item, nil
}

func (repo *repoMap) List() []item.Id {
	itemIds := make([]item.Id, len(repo.records))
	i := 0
	for id := range repo.records {
		itemIds[i] = id
		i++
	}
	return itemIds
}

func (repo *repoMap) GetAllByCompany(companyId company.Id) (items []item.Record) {
	for _, r := range repo.records {
		if r.Company() == companyId{
			items = append(items, r)
		}
	}
	return items
}
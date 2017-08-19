package factory

import "errors"
import "psychic-rat/m/item"

func GetItemRepo() item.Repo {
	return itemRepo
}


// declare that we implement Repo interface
var itemRepo item.Repo = new(repoMap)


type repoMap struct {
	records map[item.Id]item.Record
}

func (repo *repoMap) Create(i item.Record) (item.Id, error) {
	newitemId := item.Id(len(repo.records))
	repo.records[newitemId] = i
	return newitemId, nil
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

package item

import "errors"

// declare that we implement Repo interface
var _ Repo = new(RepoMap)

type RepoMap struct {
	records map[Id]Record
}

func (repo *RepoMap) Create(i Record) (Id, error) {
	newId := Id(len(repo.records))
	repo.records[newId] = i
	return newId, nil
}

func (repo *RepoMap) GetById(id Id) (Record, error) {
	item, found := repo.records[id]
	if !found {
		return Record{}, errors.New("not found")
	}
	return item, nil
}

func (repo *RepoMap) List() []Id {
	ids := make([]Id, len(repo.records))
	i := 0
	for id := range repo.records {
		ids[i] = id
		i++
	}
	return ids
}



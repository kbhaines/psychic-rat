package pledge

import (
	"errors"
	"psychic-rat/m/pubuserer"
)

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

func (repo *RepoMap) GetByUser(id pubuser.Id) []Id {
	results := make([]Id, 16)
	for _, pledge := range repo.records {
		if id == pledge.UserId {
			results = append(results, pledge.Id)
		}
	}
	return results
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


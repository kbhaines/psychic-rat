package factory

import (
	"errors"
	"psychic-rat/m/pubuser"
)

// declare that we implement Repo interface
var pubuserRepo pubuser.Repo = new(pubUserMapRepo)

func GetPubUserRepo() pubuser.Repo {
	return pubuserRepo
}

type pubUserMapRepo struct {
	records map[pubuser.Id]pubuser.Record
}

func (repo *pubUserMapRepo) Create(i pubuser.Record) (pubuser.Id, error) {
	newId := pubuser.Id(len(repo.records))
	repo.records[newId] = i
	return newId, nil
}

func (repo *pubUserMapRepo) GetById(id pubuser.Id) (pubuser.Record, error) {
	item, found := repo.records[id]
	if !found {
		return pubuser.Record{}, errors.New("not found")
	}
	return item, nil
}

func (repo *pubUserMapRepo) List() []pubuser.Id {
	ids := make([]pubuser.Id, len(repo.records))
	i := 0
	for id := range repo.records {
		ids[i] = id
		i++
	}
	return ids
}


package factory

import (
	"psychic-rat/m/pubuser"
	"psychic-rat/m/pledge"
	"errors"
)

func GetPledgeRepo() pledge.Repo {
	return pledgeRepo
}

// declare that we implement pledgeMapRepo interface
var pledgeRepo pledge.Repo = new(pledgeMapRepo)

type pledgeMapRepo struct {
	records map[pledge.Id]pledge.Record
}

func (repo *pledgeMapRepo) Create(i pledge.Record) (pledge.Id, error) {
	newId := pledge.Id(len(repo.records))
	repo.records[newId] = i
	return newId, nil
}

func (repo *pledgeMapRepo) GetById(id pledge.Id) (pledge.Record, error) {
	item, found := repo.records[id]
	if !found {
		return nil, errors.New("not found")
	}
	return item, nil
}

func (repo *pledgeMapRepo) GetByUser(id pubuser.Id) []pledge.Id {
	results := make([]pledge.Id, 16)
	for _, pledge := range repo.records {
		if id == pledge.UserId() {
			results = append(results, pledge.Id())
		}
	}
	return results
}

func (repo *pledgeMapRepo) List() []pledge.Id {
	ids := make([]pledge.Id, len(repo.records))
	i := 0
	for id := range repo.records {
		ids[i] = id
		i++
	}
	return ids
}

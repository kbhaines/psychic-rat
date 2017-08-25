package pledgerepo

import (
	"psychic-rat/repo"
	"psychic-rat/mdl/pledge"
	"errors"
	"psychic-rat/mdl/user"
	"fmt"
)

func GetPledgeRepoMapImpl() repo.Pledges {
	return pledgeRepo
}

// declare that we implement pledgeMapRepo interface
var pledgeRepo repo.Pledges = &pledgeMapRepo{make(map[pledge.Id]pledge.Record)}

type pledgeMapRepo struct {
	records map[pledge.Id]pledge.Record
}

func (repo *pledgeMapRepo) Create(i pledge.Record) error {
	if _, found := repo.records[i.Id()]; found {
		return fmt.Errorf("pledge with id %v exists", i.Id())
	}
	repo.records[i.Id()] = i
	return nil
}

func (repo *pledgeMapRepo) GetById(id pledge.Id) (pledge.Record, error) {
	item, found := repo.records[id]
	if !found {
		return nil, errors.New("not found")
	}
	return item, nil
}

func (repo *pledgeMapRepo) GetByUser(id user.Id) []pledge.Id {
	results := make([]pledge.Id, 16)
	for _, p := range repo.records {
		if id == p.UserId() {
			results = append(results, p.Id())
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

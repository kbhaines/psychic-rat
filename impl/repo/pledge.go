package repo

import (
	"errors"
	"fmt"
	"psychic-rat/mdl"
)

// declare that we implement pledgeMapRepo interface
var pledgeRepo = &pledgeRepoMap{make(map[mdl.Id]mdl.PledgeRecord)}

type pledgeRepoMap struct {
	records map[mdl.Id]mdl.PledgeRecord
}

func (repo *pledgeRepoMap) Create(i mdl.PledgeRecord) error {
	if _, found := repo.records[i.Id]; found {
		return fmt.Errorf("pledge with id %v exists", i.Id)
	}
	repo.records[i.Id] = i
	return nil
}

func (repo *pledgeRepoMap) GetById(id mdl.Id) (*mdl.PledgeRecord, error) {
	item, found := repo.records[id]
	if !found {
		return nil, errors.New("not found")
	}
	return &item, nil
}

func (repo *pledgeRepoMap) GetByUser(id mdl.Id) []mdl.Id {
	results := make([]mdl.Id, 16)
	for _, p := range repo.records {
		if id == p.UserId {
			results = append(results, p.Id)
		}
	}
	return results
}

func (repo *pledgeRepoMap) List() (pledges []mdl.PledgeRecord) {
	for _, p := range repo.records {
		pledges = append(pledges, p)
	}
	return pledges
}

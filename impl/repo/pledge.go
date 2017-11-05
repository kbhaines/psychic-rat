package repo

import (
	"errors"
	"fmt"
	"psychic-rat/mdl"
)

// declare that we implement pledgeMapRepo interface
var pledgeRepo = &pledgeRepoMap{make(map[mdl.ID]mdl.Pledge)}

type pledgeRepoMap struct {
	records map[mdl.ID]mdl.Pledge
}

func (repo *pledgeRepoMap) Create(i mdl.Pledge) error {
	if _, found := repo.records[i.Id]; found {
		return fmt.Errorf("pledge with id %v exists", i.Id)
	}
	repo.records[i.Id] = i
	return nil
}

func (repo *pledgeRepoMap) GetById(id mdl.ID) (*mdl.Pledge, error) {
	item, found := repo.records[id]
	if !found {
		return nil, errors.New("not found")
	}
	return &item, nil
}

func (repo *pledgeRepoMap) GetByUser(id mdl.ID) []mdl.ID {
	results := make([]mdl.ID, 16)
	for _, p := range repo.records {
		if id == p.UserID {
			results = append(results, p.Id)
		}
	}
	return results
}

func (repo *pledgeRepoMap) List() (pledges []mdl.Pledge) {
	for _, p := range repo.records {
		pledges = append(pledges, p)
	}
	return pledges
}

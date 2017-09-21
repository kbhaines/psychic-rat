package pledgerepo

import (
	"errors"
	"fmt"
	"psychic-rat/mdl/pledge"
	"psychic-rat/mdl/user"
	"psychic-rat/repo"
)

func GetPledgeRepoMapImpl() repo.Pledges {
	return pledgeRepo
}

// declare that we implement pledgeMapRepo interface
var pledgeRepo repo.Pledges = &pledgeMapRepo{make(map[pledge.Id]pledge.PledgeRecord)}

type pledgeMapRepo struct {
	records map[pledge.Id]pledge.PledgeRecord
}

func (repo *pledgeMapRepo) Create(i pledge.PledgeRecord) error {
	if _, found := repo.records[i.Id()]; found {
		return fmt.Errorf("pledge with id %v exists", i.Id())
	}
	repo.records[i.Id()] = i
	return nil
}

func (repo *pledgeMapRepo) GetById(id pledge.Id) (pledge.PledgeRecord, error) {
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

func (repo *pledgeMapRepo) List() (pledges []pledge.PledgeRecord) {
	for _, p := range repo.records {
		pledges = append(pledges, p)
	}
	return pledges
}

package ctr

import (
	"fmt"
	"psychic-rat/mdl/item"
	"psychic-rat/mdl/pledge"
	"psychic-rat/mdl/user"
	"psychic-rat/repo/pledgerepo"
	"time"
)

type PledgeController interface {
	AddPledge(itemId item.Id, userId user.Id) (pledge.Id, error)
	ListPledges(filter PledgeFilter) []pledge.Record
}

type PledgeFilter func(record pledge.Record) pledge.Record

var pledgeRepo = pledgerepo.GetPledgeRepoMapImpl()

var _ PledgeController = &pledgeController{}

type pledgeController struct{}

func (p *pledgeController) AddPledge(itemId item.Id, userId user.Id) (newId pledge.Id, err error) {
	_, err = itemRepo.GetById(itemId)
	if err != nil {
		return newId, fmt.Errorf("error retrieving item %v: %v", itemId, err)
	}
	_, err = userRepo.GetById(userId)
	if userId != 0 && err != nil {
		return newId, fmt.Errorf("error retrieving user %v: %v", userId, err)
	}
	newPledge := pledge.New(userId, itemId, time.Now())
	pledgeRepo.Create(newPledge)
	return newPledge.Id(), nil
}

func (p *pledgeController) ListPledges(filter PledgeFilter) (pledges []pledge.Record) {
	if filter == nil {
		filter = func(i pledge.Record) pledge.Record { return i }
	}
	ps := pledgeRepo.List()
	for _, p := range ps {
		if filtered := filter(p); filtered != nil {
			pledges = append(pledges, filtered)
		}
	}
	return pledges
}

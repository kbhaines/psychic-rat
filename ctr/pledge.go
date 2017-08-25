package ctr

import (
	"psychic-rat/mdl/user"
	"psychic-rat/mdl/pledge"
	"time"
	"fmt"
	"psychic-rat/mdl/item"
	"psychic-rat/repo/pledgerepo"
)

type PledgeController interface {
	AddPledge(itemId item.Id, userId user.Id) error
}

var pledgeRepo = pledgerepo.GetPledgeRepoMapImpl()

var _ PledgeController = &pledgeController{}

type pledgeController struct{}

func (p *pledgeController) AddPledge(itemId item.Id, userId user.Id) error {
	_, err := itemRepo.GetById(itemId)
	if err != nil {
		return fmt.Errorf("error retrieving item %v: %v", itemId, err)
	}
	_, err = userRepo.GetById(userId)
	if userId != 0 && err != nil {
		return fmt.Errorf("error retrieving user %v: %v", userId, err)
	}
	newPledge := pledge.New(userId, itemId, time.Now())
	pledgeRepo.Create(newPledge)
	return nil
}
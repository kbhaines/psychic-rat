package ctr

import (
	"psychic-rat/mdl/item"
	"psychic-rat/mdl/pubuser"
	"psychic-rat/mdl/pledge"
	"time"
	"psychic-rat/factory"
	"fmt"
)

type PledgeController interface {
	MakeAddPledgeRequest(itemId item.Id, userId pubuser.Id) AddPledgeRequest
	HandlePledgeRequest(req AddPledgeRequest) error
}

type AddPledgeRequest interface {
	ItemId() item.Id
	UserId() pubuser.Id
}

type addPledgeRequest struct {
	itemId item.Id
	userId pubuser.Id
}

func (p *addPledgeRequest) ItemId() item.Id {
	return p.itemId
}

func (p *addPledgeRequest) UserId() pubuser.Id {
	return p.userId
}

var itemRepo = factory.GetItemRepo()
var pledgeRepo = factory.GetPledgeRepo()
var userRepo = factory.GetPubUserRepo()

var _ PledgeController = &pledgeController{}

type pledgeController struct{}

func (p *pledgeController) MakeAddPledgeRequest(itemId item.Id, userId pubuser.Id) AddPledgeRequest {
	return &addPledgeRequest{itemId: itemId, userId: userId}
}

func (p *pledgeController) HandlePledgeRequest(req AddPledgeRequest) error {
	_, err := itemRepo.GetById(req.ItemId())
	if err != nil {
		return fmt.Errorf("error retrieving item %v: %v", req.ItemId(), err)
	}
	_, err = userRepo.GetById(req.UserId())
	if err != nil {
		return fmt.Errorf("error retrieving user %v: %v", req.ItemId(), err)
	}
	newPledge := pledge.New(req.UserId(), req.ItemId(), time.Now())
	pledgeRepo.Create(newPledge)
	return nil
}

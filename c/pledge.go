package c

import (
	"psychic-rat/m/item"
	"psychic-rat/m/pubuser"
	"psychic-rat/m/pledge"
	"time"
	"psychic-rat/factory"
	"fmt"
)

type NewPledgeRequest interface {
	ItemId() item.Id
	UserId() pubuser.Id
}

type newPledgeRequest struct {
	itemId item.Id
	userId pubuser.Id
}

func (p *newPledgeRequest) ItemId() item.Id {
	return p.itemId
}

func (p *newPledgeRequest) UserId() pubuser.Id {
	return p.userId
}

func MakePledgeRequest(itemId item.Id, userId pubuser.Id) NewPledgeRequest {
	return &newPledgeRequest{itemId: itemId, userId: userId}
}

var itemRepo = factory.GetItemRepo()
var pledgeRepo = factory.GetPledgeRepo()
var userRepo = factory.GetPubUserRepo()

func HandlePledgeRequest(req NewPledgeRequest) error {
	_, err := itemRepo.GetById(req.ItemId())
	if err != nil {
		return fmt.Errorf("error retrieving item %v: %v", req.ItemId(), err )
	}
	_, err = userRepo.GetById(req.UserId())
	if err != nil {
		return fmt.Errorf("error retrieving user %v: %v", req.ItemId(), err )
	}
	newPledge := pledge.New(req.UserId(), req.ItemId(), time.Now())
	pledgeRepo.Create(newPledge)
	return nil
}

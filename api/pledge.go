package api

import (
	"fmt"
	"psychic-rat/mdl/item"
	"psychic-rat/mdl/pledge"
	"psychic-rat/mdl/user"
	"psychic-rat/repo/itemrepo"
	"psychic-rat/repo/pledgerepo"
	"psychic-rat/repo/userrepo"
	"time"
)

type PledgeController interface {
	AddPledge(itemId item.Id, userId user.Id) (pledge.Id, error)
	ListPledges(filter PledgeFilter) PledgeListing
}

type PledgeListing struct {
	Pledges []PledgeElement `json:"pledges"`
}

type PledgeElement struct {
	PledgeId  pledge.Id   `json:"id"`
	Item      ItemElement `json:"item"`
	Timestamp time.Time   `json:"timestamp"`
}

type PledgeRequest struct {
	ItemId item.Id `json:"itemId"`
}

type PledgeFilter func(record pledge.Record) pledge.Record

type pledgeApiRepoImpl struct{}

func (p *pledgeApiRepoImpl) AddPledge(itemId item.Id, userId user.Id) (newId pledge.Id, err error) {
	_, err = itemrepo.GetItemRepoMapImpl().GetById(itemId)
	if err != nil {
		return newId, fmt.Errorf("error retrieving item %v: %v", itemId, err)
	}
	_, err = userrepo.GetUserRepoMapImpl().GetById(userId)
	if userId != 0 && err != nil {
		return newId, fmt.Errorf("error retrieving user %v: %v", userId, err)
	}
	newPledge := pledge.New(userId, itemId, time.Now())
	pledgerepo.GetPledgeRepoMapImpl().Create(newPledge)
	return newPledge.Id(), nil
}

func (p *pledgeApiRepoImpl) ListPledges(filter PledgeFilter) (PledgeListing, error) {
	if filter == nil {
		filter = func(i pledge.Record) pledge.Record { return i }
	}
	return PledgeListing{}, nil
}

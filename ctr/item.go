package ctr

import (
	"psychic-rat/mdl/company"
	"fmt"
	"psychic-rat/mdl/item"
	"psychic-rat/repo/itemrepo"
)

type ItemController interface {
	AddItem(make string, model string, company company.Id) error
	ListItems() []item.Record
}

type itemController struct{}

var _ ItemController = &itemController{}

var itemRepo = itemrepo.GetItemRepoMapImpl()

func (i *itemController) AddItem(make string, model string, company company.Id) error {
	item := item.New(make, model, company)
	err := checkDuplicate(item)
	if err != nil {
		return fmt.Errorf("duplicate check failed for new item %s/%s/%s: %v", make, model, company, err)
	}

	err = itemRepo.Create(item)
	if err != nil {
		return fmt.Errorf("couldn't create item %v: %v", item, err)
	}
	return nil
}

func (i *itemController) ListItems() (items []item.Record) {
	itemIds := itemRepo.List()
	for _, i := range itemIds {
		item, _ := itemRepo.GetById(i)
		items = append(items, item)
	}
	return
}

func checkDuplicate(item item.Record) error {
	itemsToCheck := itemRepo.GetAllByCompany(item.Company())
	for _, i := range itemsToCheck {
		if item.Make() == i.Make() && item.Model() == i.Model() {
			return fmt.Errorf("existing item: %v", i)
		}
	}
	return nil
}

package ctr

import (
	"psychic-rat/mdl/company"
	"fmt"
	"psychic-rat/mdl/item"
	"psychic-rat/repo/itemrepo"
)

type ItemController interface {
	AddItem(make string, model string, company company.Id) error
	GetById(id item.Id) (item.Record, error)
	ListItems(f ItemFilter) []item.Record
}

type ItemFilter func(record item.Record) item.Record

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

func (i *itemController) ListItems(filter ItemFilter) (items []item.Record) {
	if filter == nil {
		filter = func(i item.Record) item.Record { return i }
	}
	itemIds := itemRepo.List()
	for _, i := range itemIds {
		item, _ := itemRepo.GetById(i)
		if filtered := filter(item); filtered != nil {
			items = append(items, filtered)
		}
	}
	return items
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

func (i *itemController) GetById(id item.Id) (item.Record, error) {
	return itemRepo.GetById(id)
}

package ctr

import (
	"psychic-rat/mdl/company"
	"fmt"
)

type ItemController interface {
	HandleAddItemRequest(request AddItemRequest) error
	MakeAddItemRequest(make string, model string, company company.Id) AddItemRequest
}

type AddItemRequest interface {
	Make() string
	Model() string
	CompanyId() company.Id
}

type addItemRequest struct {
	make      string
	model     string
	companyId company.Id
}

func (i *addItemRequest) Make() string {
	return i.make
}

func (i *addItemRequest) Model() string {
	return i.model
}

func (i *addItemRequest) CompanyId() company.Id {
	return i.companyId
}

type itemController struct{}

var _ ItemController = &itemController{}

func (i *itemController) MakeAddItemRequest(make string, model string, company company.Id) AddItemRequest {
	return &addItemRequest{make: make, model: model, companyId: company}
}

func (i *itemController) HandleAddItemRequest(req AddItemRequest) error {
	err := checkDuplicate(req)
	if err != nil {
		fmt.Errorf("duplicate check failed for %v: %v", req, err)
	}
	return nil
}

func checkDuplicate(request AddItemRequest) error {
	return nil
}
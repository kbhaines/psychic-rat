package c

import (
	"psychic-rat/m/company"
	"fmt"
)

type NewItemRequest interface {
	Make() string
	Model() string
	CompanyId() company.Id
}

type newItemRequest struct {
	make      string
	model     string
	companyId company.Id
}

func (i *newItemRequest) Make() string {
	return i.make
}

func (i *newItemRequest) Model() string {
	return i.model
}

func (i *newItemRequest) CompanyId() company.Id {
	return i.companyId
}

func MakeItemRequest(make string, model string, company company.Id) NewItemRequest {
	return &newItemRequest{make: make, model: model, companyId: company}
}

func HandleNewItemRequest(req NewItemRequest) error {
	err := checkDuplicate(req)
	if err != nil {
		fmt.Errorf("duplicate check failed for %v: %v", req, err)
	}
	return nil
}

func checkDuplicate(request NewItemRequest) error {
	return nil
}

package api

import (
	"psychic-rat/mdl/item"
	"psychic-rat/repo/itemrepo"
)

type ItemApi interface {
	//AddItem(make string, model string, company company.Id) error
	ListItems(f ItemFilter) (ItemReport, error)
	GetById(id item.Id) (ItemElement, error)
}

type ItemReport struct {
	Items []ItemElement `json:"items"`
}

type ItemElement struct {
	Id    item.Id `json:"id"`
	Make  string  `json:"make"`
	Model string  `json:"model"`
}

type ItemFilter func(record item.Record) item.Record

type itemRepoApi struct{}

func (i *itemRepoApi) ListItems(filter ItemFilter) (ItemReport, error) {
	itemRepo := itemrepo.GetItemRepoMapImpl()
	if filter == nil {
		filter = func(i item.Record) item.Record { return i }
	}
	return ItemReport{}, nil
}

func (i *itemRepoApi) GetById(id item.Id) (ItemElement, error) {
	item, err := itemrepo.GetItemRepoMapImpl().GetById(id)
	return ItemElement{item.Id(), item.Make(), item.Model()}, err
}

//func checkDuplicate(item item.Record) error {
//	itemsToCheck := itemRepo.GetAllByCompany(item.Company())
//	for _, i := range itemsToCheck {
//		if item.Make() == i.Make() && item.Model() == i.Model() {
//			return fmt.Errorf("existing item: %v", i)
//		}
//	}
//	return nil
//}
//func createItem() {
//	const (
//		Make    = "make"
//		Model   = "model"
//		Company = "company"
//	)
//	params, ok := extractFormParams(request.Form, Make, Model, Company)
//	if ! ok {
//		fmt.Fprintf(writer, "form parameters missing: got %v", params)
//	}
//
//	err := ctr.GetController().Item().AddItem(params[Make], params[Model], company.id(params[Company]))
//	if err != nil {
//		errorResponse(writer, err)
//	}
//
//}

//func handleItemGet(writer http.ResponseWriter, request *http.Request) {
//	items := ctr.GetController().Item().ListItems()
//	fmt.Fprintf(writer, "items: %v", items)
//}

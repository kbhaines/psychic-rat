package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"psychic-rat/api/rest"
	"psychic-rat/mdl/item"
	"psychic-rat/repo/itemrepo"
)

type ItemApi interface {
	//AddItem(make string, model string, company company.Id) error
	ListItems() (ItemReport, error)
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

type itemRepoApi struct{}

////////////////////////////////////////////////////////////////////////////////
// ItemApi implementations

////////////////////////////////////////////////////////////////////////////////
// Repo api

func GetRepoItemApi() ItemApi {
	return &itemRepoApi{}
}

func (i *itemRepoApi) ListItems() (ItemReport, error) {
	repo := itemrepo.GetItemRepoMapImpl()
	items := repo.List()
	results := make([]ItemElement, len(items))
	for i, item := range items {
		results[i] = ItemElement{item.Id(), item.Make(), item.Model()}
	}
	return ItemReport{results}, nil
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

////////////////////////////////////////////////////////////////////////////////
// Restful api

type itemRestApi struct {
	url string
}

func GetRestfulItemApi(url string) ItemApi {
	return &itemRestApi{url}
}

func (a *itemRestApi) ListItems() (ItemReport, error) {
	resp, err := http.Get(a.url + rest.ItemApi + fmt.Sprintf("?company=%v", "1"))
	report := ItemReport{}
	if err != nil {
		return report, fmt.Errorf("get items failed: %v", err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return report, fmt.Errorf("error reading item response: %v", err)
	}
	return ItemsFromJson(bytes)
}

func (a *itemRestApi) GetById(id item.Id) (ItemElement, error) {
	panic("Not implemented")
}

func ItemsFromJson(bytes []byte) (ItemReport, error) {
	var items ItemReport
	if err := json.Unmarshal(bytes, &items); err != nil {
		return items, fmt.Errorf("failed to unmarshal items: %v", err)
	}
	return items, nil
}

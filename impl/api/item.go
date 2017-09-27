package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	a "psychic-rat/api"
	"psychic-rat/api/rest"
	"psychic-rat/mdl"
	"psychic-rat/repo"
)

////////////////////////////////////////////////////////////////////////////////
// ItemApi implementations

////////////////////////////////////////////////////////////////////////////////
// Repo api

type itemRepoApi struct {
	repos repo.Repos
}

func (i *itemRepoApi) ListItems() (a.ItemReport, error) {
	items := i.repos.Item.List()
	coRepo := i.repos.Company
	results := make([]a.ItemElement, len(items))
	for i, item := range items {
		company, _ := coRepo.GetById(item.CompanyId)
		results[i] = a.ItemElement{item.Id, item.Make, item.Model, company.Name}
	}
	return a.ItemReport{results}, nil
}

func (i *itemRepoApi) GetById(id mdl.Id) (a.ItemElement, error) {
	item, err := i.repos.Item.GetById(id)
	if err != nil {
		return a.ItemElement{}, err
	}
	co, err := i.repos.Company.GetById(item.CompanyId)
	if err != nil {
		return a.ItemElement{}, err
	}
	return a.ItemElement{item.Id, item.Make, item.Model, co.Name}, err
}

//func checkDuplicate(item mdl.Record) error {
//	itemsToCheck := itemRepo.GetAllByCompany(mdl.Company())
//	for _, i := range itemsToCheck {
//		if mdl.Make() == i.Make() && mdl.Model() == i.Model() {
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

func GetRestfulItemApi(url string) a.ItemApi {
	return &itemRestApi{url}
}

func (ia *itemRestApi) ListItems() (a.ItemReport, error) {
	resp, err := http.Get(ia.url + rest.ItemApi + fmt.Sprintf("?company=%v", "1"))
	report := a.ItemReport{}
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

func (a *itemRestApi) GetById(id mdl.Id) (a.ItemElement, error) {
	panic("Not implemented")
}

func ItemsFromJson(bytes []byte) (a.ItemReport, error) {
	var items a.ItemReport
	if err := json.Unmarshal(bytes, &items); err != nil {
		return items, fmt.Errorf("failed to unmarshal items: %v", err)
	}
	return items, nil
}

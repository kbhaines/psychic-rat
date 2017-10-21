package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"psychic-rat/api/rest"
	"psychic-rat/mdl"
	"psychic-rat/repo"
	"psychic-rat/types"
)

////////////////////////////////////////////////////////////////////////////////
// ItemApi implementations

////////////////////////////////////////////////////////////////////////////////
// Repo api

func NewItemApi(repos repo.Repos) *itemRepoApi {
	return &itemRepoApi{repos: repos}
}

type itemRepoApi struct {
	repos repo.Repos
}

func (i *itemRepoApi) ListItems() (types.ItemReport, error) {
	items := i.repos.Item.List()
	coRepo := i.repos.Company
	results := make([]types.ItemElement, len(items))
	for i, item := range items {
		company, _ := coRepo.GetById(item.CompanyId)
		results[i] = types.ItemElement{item.Id, item.Make, item.Model, company.Name}
	}
	return types.ItemReport{results}, nil
}

func (i *itemRepoApi) GetById(id mdl.Id) (types.ItemElement, error) {
	item, err := i.repos.Item.GetById(id)
	if err != nil {
		return types.ItemElement{}, err
	}
	co, err := i.repos.Company.GetById(item.CompanyId)
	if err != nil {
		return types.ItemElement{}, err
	}
	return types.ItemElement{item.Id, item.Make, item.Model, co.Name}, err
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

func GetRestfulItemApi(url string) *itemRestApi {
	return &itemRestApi{url}
}

func (ia *itemRestApi) ListItems() (types.ItemReport, error) {
	resp, err := http.Get(ia.url + rest.ItemApi + fmt.Sprintf("?company=%v", "1"))
	report := types.ItemReport{}
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

func (a *itemRestApi) GetById(id mdl.Id) (types.ItemElement, error) {
	panic("Not implemented")
}

func ItemsFromJson(bytes []byte) (types.ItemReport, error) {
	var items types.ItemReport
	if err := json.Unmarshal(bytes, &items); err != nil {
		return items, fmt.Errorf("failed to unmarshal items: %v", err)
	}
	return items, nil
}

package api

import (
	"fmt"
	"psychic-rat/ctr"
	"net/http"
	"psychic-rat/mdl/company"
	"psychic-rat/mdl/item"
	"encoding/json"
)

func ItemHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		unsupportedMethod(writer)
		return
	}
	if err := request.ParseForm(); err != nil {
		unableToParseForm(err, writer)
		return
	}

	companyId := request.Form.Get("company")
	if companyId == "" {
		errorResponse(writer, fmt.Errorf("no company specified"))
		return
	}
	json, err := json.Marshal(getItemsForCompany(company.Id(companyId)))
	if err != nil {
		logInternalError(writer, err)
		return
	}
	fmt.Fprintf(writer, "%s", json)
}

type items struct {
	Items []item.Record `json:"items""`
}

type itemElement struct {
	IId    item.Id `json:"id"`
	MMake  string `json:"make"`
	MModel string `json:"model"`
}

func (i *itemElement) Id() item.Id             { return i.IId }
func (i *itemElement) Make() string            { return i.MMake }
func (i *itemElement) Model() string           { return i.MModel }
func (i *itemElement) Company() (c company.Id) { return }

func getItemsForCompany(companyId company.Id) items {
	is := ctr.GetController().Item().ListItems(func(i item.Record) item.Record {
		if companyId == i.Company() {
			return &itemElement{i.Id(), i.Make(), i.Model()}
		}
		return nil
	})
	return items{is}
}

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

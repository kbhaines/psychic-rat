package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"psychic-rat/ctr"
	"psychic-rat/mdl/company"
	"psychic-rat/mdl/item"
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

type itemReport struct {
	Items []itemElement `json:"items"`
}

type itemElement struct {
	IId    item.Id `json:"id"`
	MMake  string  `json:"make"`
	MModel string  `json:"model"`
}

func (i *itemElement) Id() item.Id             { return i.IId }
func (i *itemElement) Make() string            { return i.MMake }
func (i *itemElement) Model() string           { return i.MModel }
func (i *itemElement) Company() (c company.Id) { return }

func getItemsForCompany(companyId company.Id) itemReport {
	is := ctr.GetController().Item().ListItems(func(i item.Record) item.Record {
		if companyId == i.Company() {
			return &itemElement{i.Id(), i.Make(), i.Model()}
		}
		return nil
	})
	report := itemReport{}
	for i, v := range is {
		report.Items[i] = itemElement{v.Id(), v.Make(), v.Model()}
	}
	return report
}

func ItemsFromJson(bytes []byte) ([]item.Record, error) {
	var items itemReport
	if err := json.Unmarshal(bytes, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %v", err)
	}
	results := make([]item.Record, len(items.Items))
	for i, v := range items.Items {
		results[i] = &itemElement{v.Id(), v.Make(), v.Model()}
	}
	return results, nil
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

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

type ItemReport struct {
	Items []ItemElement `json:"items"`
}

type ItemElement struct {
	Id    item.Id `json:"id"`
	Make  string  `json:"make"`
	Model string  `json:"model"`
}

func ifElse(b bool, t, f interface{}) interface{} {
	if b {
		return t
	}
	return f
}

func getItemsForCompany(companyId company.Id) ItemReport {
	is := ctr.GetController().Item().ListItems(func(i item.Record) item.Record {
		if companyId == i.Company() {
			return i
		} else {
			return nil
		}
	})
	report := ItemReport{make([]ItemElement, len(is))}
	for i, v := range is {
		report.Items[i] = ItemElement{v.Id(), v.Make(), v.Model()}
	}
	return report
}

func ItemsFromJson(bytes []byte) (ItemReport, error) {
	var items ItemReport
	if err := json.Unmarshal(bytes, &items); err != nil {
		return items, fmt.Errorf("failed to unmarshal items: %v", err)
	}
	return items, nil
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

package api

import (
	"fmt"
	"psychic-rat/ctr"
	"net/http"
	"psychic-rat/mdl/company"
)

func itemHandler(writer http.ResponseWriter, request *http.Request) {
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
	//items := getItemsForCompany(companyId)
}

type itemElements struct {
	Id string `json:"id"`
	Make string `json:"make"`
	Model string `json:"model"`
}

func getItemsForCompany(companyId company.Id) (items []itemElements) {
	return items
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
//	err := ctr.GetController().Item().AddItem(params[Make], params[Model], company.Id(params[Company]))
//	if err != nil {
//		errorResponse(writer, err)
//	}
//
//}

func handleItemGet(writer http.ResponseWriter, request *http.Request) {
	items := ctr.GetController().Item().ListItems()
	fmt.Fprintf(writer, "items: %v", items)
}

package api

import (
	"fmt"
	"psychic-rat/ctr"
	"psychic-rat/mdl/company"
	"net/http"
)

func itemHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		handleItemGet(writer, request)
		return
	}
	if request.Method != http.MethodPost {
		unsupportedMethod(writer)
		return
	}
	if err := request.ParseForm(); err != nil {
		unableToParseForm(err, writer)
		return
	}

	const (
		Make    = "make"
		Model   = "model"
		Company = "company"
	)
	params, ok := extractFormParams(request.Form, Make, Model, Company)
	if ! ok {
		fmt.Fprintf(writer, "form parameters missing: got %v", params)
	}

	err := ctr.GetController().Item().AddItem(params[Make], params[Model], company.Id(params[Company]))
	if err != nil {
		errorResponse(writer, err)
	}

}

func handleItemGet(writer http.ResponseWriter, request *http.Request) {
	items := ctr.GetController().Item().ListItems()
	fmt.Fprintf(writer, "items: %v", items)
}

package main

import (
	"net/http"
	"fmt"
	"log"
	"net/url"
	"psychic-rat/mdl/item"
	"psychic-rat/mdl/user"
	"psychic-rat/ctr"
	"psychic-rat/mdl/company"
)

type MethodHandler func(http.ResponseWriter, *http.Request)

type HandlerMap struct {
	PathExpr string
	Method   string
	Handler  MethodHandler
}

func main() {
	http.HandleFunc("/api/v1/pledge", pledgeHandler)
	http.HandleFunc("/api/v1/item", itemHandler)
	http.ListenAndServe("localhost:8080", nil)
}

func pledgeHandler(writer http.ResponseWriter, request *http.Request) {

	handlerMaps := map[string]MethodHandler{
		http.MethodPost: handlePledgePost,
		http.MethodGet:  handleGet,
	}

	fmt.Printf("%v", request)
	v, ok := handlerMaps[request.Method]
	if ! ok {
		unsupportedMethod(writer)
		return
	}
	v(writer, request)
}

func unsupportedMethod(writer http.ResponseWriter) {
	fmt.Fprintf(writer, "unsupported method")
}

func handlePledgePost(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		unableToParseForm(err, writer)
		return
	}
	itemId, userId, err := parsePledgePost(request.Form)
	if err != nil {
		fmt.Fprintf(writer, "error: %v", err)
		return
	}

	err = ctr.GetController().Pledge().AddPledge(itemId, userId)
	if err != nil {
		fmt.Fprintf(writer, "error: %v", err)
	}
}

func unableToParseForm(err error, writer http.ResponseWriter) {
	fmt.Fprintf(writer, "error in form data")
	log.Print(err)
}

func parsePledgePost(values url.Values) (itemId item.Id, userId user.Id, err error) {
	const (
		Item = "item"
	)

	params, ok := extractFormParams(values, Item)
	if ! ok {
		return itemId, userId, fmt.Errorf("missing values, only got %v", params)
	}

	return item.Id(params[Item]), user.Id(0), nil
}

func extractFormParams(values url.Values, params ...string) (results map[string]string, gotAllParams bool) {
	results = make(map[string]string)
	gotAllParams = true
	for _, p := range (params) {
		v, ok := values[p]
		if ! ok {
			gotAllParams = false
			continue
		}
		results[p] = v[0]
	}
	return results, gotAllParams
}

func handleGet(writer http.ResponseWriter, request *http.Request) {
	const message = "Welcome"
	fmt.Fprintf(writer, message)
}

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

func errorResponse(writer http.ResponseWriter, err error) {
	fmt.Fprintf(writer, "error: %v", err)
}

func handleItemGet(writer http.ResponseWriter, request *http.Request) {
	items := ctr.GetController().Item().ListItems()
	fmt.Fprintf(writer, "items: %v", items)
}

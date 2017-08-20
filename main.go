package main

import (
	"net/http"
	"fmt"
	"log"
	"net/url"
	"psychic-rat/mdl/item"
	"psychic-rat/mdl/pubuser"
	"psychic-rat/ctr"
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
	pledge, err := parsePledgePost(request.Form)
	if err != nil {
		fmt.Fprintf(writer, "error: %v", err)
		return
	}

	err = ctr.GetController().Pledge().HandlePledgeRequest(pledge)
	if err != nil {
		fmt.Fprintf(writer, "error: %v", err)
	}
}

func unableToParseForm(err error, writer http.ResponseWriter) {
	fmt.Fprintf(writer, "error in form data")
	log.Print(err)
}

func parsePledgePost(values url.Values) (ctr.NewPledgeRequest, error) {
	const (
		Item = "item"
	)

	params, ok := extractFormParams(values, Item)
	if ! ok {
		return nil, fmt.Errorf("missing values, only got %v", params)
	}

	return ctr.GetController().Pledge().MakePledgeRequest(item.Id(params[Item]), pubuser.Id(0)), nil
}

func extractFormParams(values url.Values, params ...string) (map[string]string, bool) {
	results := map[string]string{}
	resultOk := true
	for _, p := range (params) {
		v, ok := values[p]
		if ! ok {
			resultOk = false
			continue
		}
		results[p] = v[0]
	}
	return results, resultOk
}

func handleGet(writer http.ResponseWriter, request *http.Request) {
	const message = "Welcome"
	fmt.Fprintf(writer, message)
}

func itemHandler(writer http.ResponseWriter, request *http.Request) {
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

	ctr.GetController()

}

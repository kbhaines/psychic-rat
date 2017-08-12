package main

import (
	"net/http"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"github.com/satori/go.uuid"
	ent "psychic-rat/pledge"
)

type MethodHandler func(http.ResponseWriter, *http.Request)

type HandlerMap struct {
	PathExpr string
	Method   string
	Handler  MethodHandler
}

type PledgeRepo interface {
	Save(pledge ent.Record) (bool)
	GetByUserId(id ent.UserId) []ent.Record
}

func main() {
	http.HandleFunc("/api/v1/pledge", pledgeHandler)
	http.ListenAndServe("localhost:8080", nil)
}

func pledgeHandler(writer http.ResponseWriter, request *http.Request) {

	handlerMaps := map[string]MethodHandler{
		http.MethodPost: handlePost,
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

func handlePost(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		unableToParseForm(err, writer)
		return
	}
	pledge, err := parseRequest(request.Form)
	if err != nil {
		fmt.Fprintf(writer, "error: %v", err)
		return
	}

	//savePledge()
	fmt.Fprintf(writer, "%v", pledge)
}

func handleGet(writer http.ResponseWriter, request *http.Request) {
	const message = "Welcome"
	fmt.Fprintf(writer, message)
}

func parseRequest(values url.Values) (ent.Record, error) {
	const (
		Item    = "item"
		Company = "company"
		Value   = "value"
		Email   = "email"
		Country = "country"
	)

	params, ok := extractFormParams(values, Item, Company, Value, Email)
	if ! ok {
		return ent.Record{}, fmt.Errorf("missing values, only got %v", params)
	}

	value, err := strconv.Atoi(params[Value])
	if err != nil {
		return ent.Record{}, err
	}
	newId := uuid.NewV4().String()
	return ent.Record{ent.Id(newId), params[Email], ent.ItemId(params[Item]), ent.CompanyId(params[Company]), value}, nil
}

func extractFormParams(values url.Values, params ...string) (map[string]string, bool) {
	results := map[string]string{}
	for _, p := range (params) {
		v, ok := values[p]
		if ! ok {
			return results, false
		}
		results[p] = v[0]
	}
	return results, true
}

func unableToParseForm(err error, writer http.ResponseWriter) {
	fmt.Fprintf(writer, "error in form data")
	log.Print(err)
}

func unsupportedMethod(writer http.ResponseWriter) {
	fmt.Fprintf(writer, "unsupported method")
}

package main

import (
	"net/http"
	"fmt"
	"log"
	"net/url"
	"strconv"
)

func main() {

	http.HandleFunc("/api/v1/pledge", handler)
	http.ListenAndServe("localhost:8080", nil)

}

type MethodHandler func(http.ResponseWriter, *http.Request)

func handler(writer http.ResponseWriter, request *http.Request) {

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
	fmt.Fprintf(writer, "%v", pledge)
}

func handleGet(writer http.ResponseWriter, request *http.Request) {
	const message = "Welcome"
	fmt.Fprintf(writer, message)
}

type CompanyId string

type ItemId string

type Pledge struct {
	Email     string
	ItemId    ItemId
	CompanyId CompanyId
	Value     int
}

func parseRequest(values url.Values) (Pledge, error) {
	const (
		Type    = "type"
		Company = "company"
		Value   = "value"
		Email   = "email"
	)

	params, ok := extractFormParams(values, Type, Company, Value, Email)
	if ! ok {
		return Pledge{}, fmt.Errorf("missing values, only got %v", params)
	}

	value, err := strconv.Atoi(params[Value])
	if err != nil {
		return Pledge{}, err
	}
	return Pledge{params[Email], ItemId(params[Type]), CompanyId(params[Company]), value}, nil
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

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
func handler(writer http.ResponseWriter, request *http.Request) {

	fmt.Printf("%v", request)
	if request.Method != http.MethodPost {
		unsupportedMethod(writer)
		return
	}
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

type Pledge struct {
	Type    string
	Company string
	value   int
}

func parseRequest(values url.Values) (Pledge, error) {
	params, ok := checkFormParametersExist(values, "type", "company", "value")
	if ! ok {
		return Pledge{}, fmt.Errorf("missing values, only got %v", params)
	}
	value, err := strconv.Atoi(params["value"])
	if err != nil {
		return Pledge{}, err
	}
	return Pledge{params["type"], params["company"], value}, nil
}

func checkFormParametersExist(values url.Values, params ...string) (map[string]string, bool) {
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

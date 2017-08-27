package api

import (
	"fmt"
	"net/http"
	"log"
	"net/url"
	"psychic-rat/mdl/user"
)

type MethodHandler func(http.ResponseWriter, *http.Request)

type HandlerMap struct {
	PathExpr string
	Method   string
	Handler  MethodHandler
}


func unsupportedMethod(writer http.ResponseWriter) {
	fmt.Fprintf(writer, "unsupported method")
}

func unableToParseForm(err error, writer http.ResponseWriter) {
	fmt.Fprintf(writer, "error in form data")
	log.Print(err)
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

func errorResponse(writer http.ResponseWriter, err error) {
	fmt.Fprintf(writer, "error: %v", err)
}

func logInternalError(writer http.ResponseWriter, err error) {
	fmt.Fprintf(writer, "internal error; contact developer: %v", err)
}

func getCurrentUserId() user.Id {
	return user.Id(0)
}
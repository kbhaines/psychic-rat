package main

import (
	"net/http"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"psychic-rat/m/item"
	"psychic-rat/m/pubuser"
	"psychic-rat/factory"
	"psychic-rat/m/pledge"
	"time"
)

type MethodHandler func(http.ResponseWriter, *http.Request)

type HandlerMap struct {
	PathExpr string
	Method   string
	Handler  MethodHandler
}

type PledgePost struct {
	itemId item.Id
	userId pubuser.Id
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

	err = validateRequest(pledge)
	if err != nil {
		fmt.Fprintf(writer, "error: %v", err)
	}
}

var itemRepo = factory.GetItemRepo()
var pledgeRepo = factory.GetPledgeRepo()
var userRepo = factory.GetPubUserRepo()

func validateRequest(post *PledgePost) error {
	_, err := itemRepo.GetById(post.itemId)
	if err != nil {
		return err
	}
	_, err = userRepo.GetById(post.userId)
	if err != nil {
		return err
	}

	newPledge := pledge.New(post.userId, post.itemId, time.Now())
	pledgeRepo.Create(newPledge)
	return nil
}

func handleGet(writer http.ResponseWriter, request *http.Request) {
	const message = "Welcome"
	fmt.Fprintf(writer, message)
}

func parseRequest(values url.Values) (*PledgePost, error) {
	const (
		Item = "item"
	)

	params, ok := extractFormParams(values, Item)
	if ! ok {
		return nil, fmt.Errorf("missing values, only got %v", params)
	}

	id, err := strconv.Atoi(params[Item])
	if err != nil {
		return nil, fmt.Errorf("illegal value for item (%v)", params[Item])
	}

	return &PledgePost{item.Id(id), pubuser.Id(0)}, nil
}

func extractFormParams(values url.Values, params ...string) (map[string]string, bool) {
	results := map[string]string{}
	resultOk := true
	for _, p := range (params) {
		v, ok := values[p]
		if ! ok {
			resultOk = false
		}
		results[p] = v[0]
	}
	return results, resultOk
}

func unableToParseForm(err error, writer http.ResponseWriter) {
	fmt.Fprintf(writer, "error in form data")
	log.Print(err)
}

func unsupportedMethod(writer http.ResponseWriter) {
	fmt.Fprintf(writer, "unsupported method")
}

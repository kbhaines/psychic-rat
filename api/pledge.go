package api

import (
	"net/http"
	"fmt"
	"psychic-rat/ctr"
	"net/url"
	"psychic-rat/mdl/item"
	"psychic-rat/mdl/user"
)

type MethodHandler func(http.ResponseWriter, *http.Request)

func pledgeHandler(writer http.ResponseWriter, request *http.Request) {

	handlerMaps := map[string]MethodHandler{
		http.MethodPost: handlePledgePost,
	}

	fmt.Printf("%v", request)
	v, ok := handlerMaps[request.Method]
	if ! ok {
		unsupportedMethod(writer)
		return
	}
	v(writer, request)
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

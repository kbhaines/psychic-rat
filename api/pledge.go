package api

import (
	"net/http"
	"fmt"
	"psychic-rat/ctr"
	"net/url"
	"psychic-rat/mdl/item"
	"psychic-rat/mdl/user"
	"encoding/json"
	"io/ioutil"
	"psychic-rat/mdl/pledge"
	"time"
)

type MethodHandler func(http.ResponseWriter, *http.Request)

func PledgeHandler(writer http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case http.MethodPost:
		doPostRequest(writer, request)

	case http.MethodGet:
		doGetRequest(writer, request)

	default:
		unsupportedMethod(writer)
	}
}

type pledgeRequest struct {
	ItemId item.Id `json:"itemId"`
	userId user.Id
}

type pledgeListing struct {
	Pledges []pledgeElement `json:"pledges"`
}

type pledgeElement struct {
	PledgeId  pledge.Id `json:"id"`
	Item      item.Record `json:"item"`
	Timestamp time.Time `json:"timestamp"`
}

func doPostRequest(writer http.ResponseWriter, request *http.Request) {

	pledge := pledgeRequest{userId: 0}

	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		logInternalError(writer, err)
		return
	}
	err = json.Unmarshal(body, &pledge)
	if err != nil {
		logInternalError(writer, err)
	}
	if err = ctr.GetController().Pledge().AddPledge(pledge.ItemId, pledge.userId); err != nil {
		logInternalError(writer, err)
	}

	fmt.Fprintf(writer, "added")
}

func doGetRequest(writer http.ResponseWriter, request *http.Request) {

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

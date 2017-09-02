package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"psychic-rat/api"
	"psychic-rat/mdl/company"
	"psychic-rat/mdl/pledge"
	"psychic-rat/mdl/user"
	"sort"
)

var companyApi = api.GetRepoCompanyApi()
var itemApi = api.GetRepoItemApi()

func CompanyHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		unsupportedMethod(writer)
		return
	}
	companies, err := companyApi.GetCompanies()
	if err != nil {
		errorResponse(writer, err)
		return
	}
	ToJson(writer, companies)

}

//// ITEMS

func ItemHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		unsupportedMethod(writer)
		return
	}
	if err := request.ParseForm(); err != nil {
		unableToParseForm(err, writer)
		return
	}

	companyId := request.Form.Get("company")
	if companyId == "" {
		errorResponse(writer, fmt.Errorf("no company specified"))
		return
	}
	json, err := json.Marshal(getItemsForCompany(company.Id(companyId)))
	if err != nil {
		logInternalError(writer, err)
		return
	}
	fmt.Fprintf(writer, "%s", json)
}

func ifElse(b bool, t, f interface{}) interface{} {
	if b {
		return t
	}
	return f
}

func getItemsForCompany(companyId company.Id) api.ItemReport {
	items, err := itemApi.ListItems()
	if err != nil {
		panic(fmt.Sprintf("unable to get items", err))
	}
	return items
}

// PLEDGES

func PledgeHandler(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		handlePost(writer, request)

	case http.MethodGet:
		handleGet(writer, request)

	default:
		unsupportedMethod(writer)
	}
}

func handlePost(writer http.ResponseWriter, request *http.Request) {
	pledge := api.PledgeRequest{}

	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		logInternalError(writer, err)
		return
	}
	err = json.Unmarshal(body, &pledge)
	if err != nil {
		logInternalError(writer, err)
		return
	}
	userId := getCurrentUserId()

	//pledgeId, err := 0, nil

	if err != nil {
		logInternalError(writer, err)
		return
	}
	writeUserPledges(writer, userId)
}

func writeUserPledges(writer http.ResponseWriter, userId user.Id) {
	userPledges := getUserPledges(userId)
	log.Printf("pledges: %v", userPledges.Pledges)
	json, err := json.Marshal(userPledges)
	if err != nil {
		logInternalError(writer, err)
		return
	}
	fmt.Fprintf(writer, "%s", json)
}

func returnIfElse(b bool, ifTrue, ifFalse interface{}) interface{} {
	if b {
		return ifTrue
	} else {
		return ifFalse
	}
}

func getUserPledges(id user.Id) api.PledgeListing {
	ps := make([]pledge.Record, 0)
	sort.Sort(pledge.ByTimeStamp(ps))
	ps2 := make([]api.PledgeElement, len(ps))
	for _, p := range ps {
		var err error
		if err != nil {
			panic(fmt.Sprintf("data inconsistency error %v. item %v for pledge %v does not exist", err, p.ItemId(), p.Id()))
		}
		//ps2[i] = api.PledgeElement{p.Id(), api.ItemElement{item.Id(), item.Make(), item.Model()}, p.TimeStamp()}
	}
	return api.PledgeListing{ps2}
}

func handleGet(writer http.ResponseWriter, request *http.Request) {
	userId := getCurrentUserId()
	writeUserPledges(writer, userId)
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
	for _, p := range params {
		v, ok := values[p]
		if !ok {
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

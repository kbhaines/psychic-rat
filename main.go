package main

import (
	"net/http"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"github.com/satori/go.uuid"
	"time"
)

type CompanyId string

type ItemId string

type CountryId string

type UserId string

type PledgeId string

type Pledge struct {
	Id        PledgeId
	UserId    UserId
	ItemId    ItemId
	Timestamp time.Time
}

type PublicUser struct {
	Id        UserId
	Country   CountryId
	FirstName string
}

type PrivateUser struct {
	Id       UserId
	Email    string
	FullName string
	AuthMethod string
	AuthSecret string
}

type MethodHandler func(http.ResponseWriter, *http.Request)

type HandlerMap struct {
	PathExpr string
	Method   string
	Handler  MethodHandler
}

type PledgeRepo interface {
	Save(pledge Pledge) (bool)
	GetByUserId(id UserId) []Pledge
}

type PublicUserRepo interface {
	CreateUser(user PublicUser) (bool)
	GetById(id UserId) (PublicUser)
}

type PrivateUserRepo interface {
	CreateUser(user PrivateUser) (bool)
	GetByEmail(email string) (PublicUser)
	GetById(id UserId) (PublicUser)
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

func parseRequest(values url.Values) (Pledge, error) {
	const (
		Item    = "item"
		Company = "company"
		Value   = "value"
		Email   = "email"
		Country = "country"
	)

	params, ok := extractFormParams(values, Item, Company, Value, Email)
	if ! ok {
		return Pledge{}, fmt.Errorf("missing values, only got %v", params)
	}

	value, err := strconv.Atoi(params[Value])
	if err != nil {
		return Pledge{}, err
	}
	newId := uuid.NewV4().String()
	return Pledge{PledgeId(newId), params[Email], ItemId(params[Item]), CompanyId(params[Company]), value}, nil
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

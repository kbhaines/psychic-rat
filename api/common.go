package api

import (
	"net/http"
)

const (
	ApiRoot    = "/api/v1"
	CompanyApi = ApiRoot + "/company"
	ItemApi    = ApiRoot + "/item"
	PledgeApi  = ApiRoot + "/pledge"
)

type UriHandler struct {
	Uri     string
	Handler http.HandlerFunc
}

var UriHandlers = []UriHandler{
	{CompanyApi, CompanyHandler},
	{ItemApi, ItemHandler},
	{PledgeApi, PledgeHandler},
}

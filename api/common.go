package api

import (
	"net/http"
	"io"
	"encoding/json"
	"fmt"
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


func ToJson(writer io.Writer, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("unable to convert %T (%v)to json", v, v))
	}
	fmt.Fprintf(writer, "%s", js)
}


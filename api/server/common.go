package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	fmt.Fprintf(writer, "%s", ToJsonString(v))
}

func ToJsonString(v interface{}) string {
	js, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("unable to convert %T (%v)to json", v, v))
	}
	return string(js)

}

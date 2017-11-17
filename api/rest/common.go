package rest

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	ApiRoot    = "/api/v1"
	CompanyApi = ApiRoot + "/company"
	ItemApi    = ApiRoot + "/item"
	PledgeApi  = ApiRoot + "/pledge"

	HomePage   = "/"
	SignInPage = "/signin"
	PledgePage = "/pledge"
	ThanksPage = "/thanks"
	NewItem    = "/newitem"
)

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

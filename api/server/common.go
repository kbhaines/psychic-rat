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

package api

import (
	"net/http"
	"psychic-rat/ctr"
	"encoding/json"
	"fmt"
	"psychic-rat/mdl/company"
)

type companyElement struct {
	Id   company.Id  `json:"id"`
	Name string      `json:"name"`
}

func CompanyHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		unsupportedMethod(writer)
		return
	}
	json, err := json.Marshal(getCompanies())
	if err != nil {
		logInternalError(writer, err)
		return
	}
	fmt.Fprintf(writer, "%s", json)
}

func getCompanies() (response []companyElement) {
	companies := ctr.GetController().Company().GetCompanies()
	for _, c := range (companies) {
		rec := companyElement{c.Id(), c.Name()}
		response = append(response, rec)
	}
	return response
}

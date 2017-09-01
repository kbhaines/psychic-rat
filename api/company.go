package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"psychic-rat/ctr"
	"psychic-rat/mdl/company"
)

type CompanyListing struct {
	Companies []CompanyElement
}

type CompanyElement struct {
	Id   company.Id `json:"id"`
	Name string     `json:"name"`
}

func CompanyHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		unsupportedMethod(writer)
		return
	}
	ToJson(writer, getCompanies())
}

func getCompanies() CompanyListing {
	companies := ctr.GetController().Company().GetCompanies()
	resp := CompanyListing{make([]CompanyElement, len(companies))}
	for i, c := range companies {
		resp.Companies[i] = CompanyElement{c.Id(), c.Name()}
	}
	return resp
}

func CompaniesFromJson(bytes []byte) (CompanyListing, error) {
	companies := CompanyListing{}
	if err := json.Unmarshal(bytes, &companies); err != nil {
		return companies, fmt.Errorf("failed to unmarshal companies: %v", err)
	}
	return companies, nil
}

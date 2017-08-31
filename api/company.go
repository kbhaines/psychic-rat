package api

import (
	"net/http"
	"psychic-rat/ctr"
	"encoding/json"
	"fmt"
	"psychic-rat/mdl/company"
)

type companyElement struct {
	Cid   company.Id  `json:"id"`
	Cname string      `json:"name"`
}

func (c *companyElement) Id() company.Id { return c.Cid }
func (c *companyElement) Name() string   { return c.Cname }

func CompanyHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		unsupportedMethod(writer)
		return
	}
	ToJson(writer, getCompanies())
}

func getCompanies() (response []companyElement) {
	companies := ctr.GetController().Company().GetCompanies()
	for _, c := range companies {
		rec := companyElement{c.Id(), c.Name()}
		response = append(response, rec)
	}
	return response
}

func CompaniesFromJson(bytes []byte) ([]company.Record, error) {
	companies := make([]companyElement, 1)
	if err := json.Unmarshal(bytes, &companies); err != nil {
		return nil, fmt.Errorf("failed to unmarshal companies: %v", err)
	}
	results := make([]company.Record, len(companies))
	for i := range companies {
		results[i] = &companies[i]
	}
	return results, nil
}


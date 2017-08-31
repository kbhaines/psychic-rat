package api

import (
	"net/http"
	"psychic-rat/ctr"
	"encoding/json"
	"fmt"
	"psychic-rat/mdl/company"
	"io"
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
	CompaniesToJson(writer)
}

func CompaniesToJson(writer io.Writer) {
	js, err := json.Marshal(getCompanies())
	if err != nil {
		panic("unable to convert companies to json")
	}
	fmt.Fprintf(writer, "%s", js)
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

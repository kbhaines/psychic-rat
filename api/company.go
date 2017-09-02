package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"psychic-rat/api/rest"
	"psychic-rat/mdl/company"
	"psychic-rat/repo/companyrepo"
)

type CompanyApi interface {
	GetCompanies() (CompanyListing, error)
}

type CompanyListing struct {
	Companies []CompanyElement `json:"companies"`
}

type CompanyElement struct {
	Id   company.Id `json:"id"`
	Name string     `json:"name"`
}

type CompanyId string

////////////////////////////////////////////////////////////////////////////////
// CompanyApi implementations

////////////////////////////////////////////////////////////////////////////////
// Repo implementation

func GetRepoCompanyApi() CompanyApi {
	return &companyApiRepoImpl{}
}

type companyApiRepoImpl struct{}

func (c *companyApiRepoImpl) GetCompanies() (CompanyListing, error) {
	companies := companyrepo.GetCompanyRepoMapImpl().GetCompanies()
	results := CompanyListing{make([]CompanyElement, len(companies))}
	for i, c := range companies {
		results.Companies[i] = CompanyElement{c.Id(), c.Name()}
	}
	return results, nil
}

////////////////////////////////////////////////////////////////////////////////
// RESTful implementation

func GetRestfulCompanyApi(baseUrl string) CompanyApi {
	return &restfulCompanyApi{baseUrl}
}

type restfulCompanyApi struct {
	url string
}

func (r *restfulCompanyApi) GetCompanies() (CompanyListing, error) {
	resp, err := http.Get(r.url + rest.CompanyApi)
	if err != nil {
		return CompanyListing{}, fmt.Errorf("get companies failed: %v", err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return CompanyListing{}, fmt.Errorf("error reading company response: %v", err)
	}
	return companiesFromJson(bytes)
}

func companiesFromJson(bytes []byte) (CompanyListing, error) {
	companies := CompanyListing{}
	if err := json.Unmarshal(bytes, &companies); err != nil {
		return companies, fmt.Errorf("failed to unmarshal companies: %v", err)
	}
	return companies, nil
}

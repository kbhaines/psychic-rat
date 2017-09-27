package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	a "psychic-rat/api"
	"psychic-rat/api/rest"
	"psychic-rat/mdl"
	"psychic-rat/repo"
)

////////////////////////////////////////////////////////////////////////////////
// CompanyApi implementations

////////////////////////////////////////////////////////////////////////////////
// Repo implementation

type companyApiRepoImpl struct {
	repos repo.Repos
}

func (c *companyApiRepoImpl) GetCompanies() (a.CompanyListing, error) {
	companies := c.repos.Company.GetCompanies()
	results := a.CompanyListing{make([]a.CompanyElement, len(companies))}
	for i, co := range companies {
		results.Companies[i] = a.CompanyElement{co.Id, co.Name}
	}
	return results, nil
}

func (c *companyApiRepoImpl) GetById(id mdl.Id) (a.CompanyElement, error) {
	co, err := c.repos.Company.GetById(id)
	if err != nil {
		return a.CompanyElement{}, err
	}
	return a.CompanyElement{id, co.Name}, nil
}

////////////////////////////////////////////////////////////////////////////////
// RESTful implementation

func GetRestfulCompanyApi(baseUrl string) a.CompanyApi {
	return &restfulCompanyApi{baseUrl}
}

type restfulCompanyApi struct {
	url string
}

func (r *restfulCompanyApi) GetCompanies() (a.CompanyListing, error) {
	resp, err := http.Get(r.url + rest.CompanyApi)
	if err != nil {
		return a.CompanyListing{}, fmt.Errorf("get companies failed: %v", err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return a.CompanyListing{}, fmt.Errorf("error reading company response: %v", err)
	}
	return companiesFromJson(bytes)
}

func companiesFromJson(bytes []byte) (a.CompanyListing, error) {
	companies := a.CompanyListing{}
	if err := json.Unmarshal(bytes, &companies); err != nil {
		return companies, fmt.Errorf("failed to unmarshal companies: %v", err)
	}
	return companies, nil
}

func (r *restfulCompanyApi) GetById(id mdl.Id) (a.CompanyElement, error) {
	resp, err := http.Get(r.url + rest.CompanyApi + fmt.Sprintf("?company=%v", id))
	if err != nil {
		return a.CompanyElement{}, fmt.Errorf("could not retrieve company %v : %v", id, err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return a.CompanyElement{}, fmt.Errorf("error reading company element from response: %v", err)
	}
	return companyFromJson(bytes)
}

func companyFromJson(bytes []byte) (a.CompanyElement, error) {
	co := a.CompanyElement{}
	if err := json.Unmarshal(bytes, &co); err != nil {
		return co, fmt.Errorf("failed to unmarshal company: %v", err)
	}
	return co, nil
}

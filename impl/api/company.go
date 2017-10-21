package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"psychic-rat/api/rest"
	"psychic-rat/mdl"
	"psychic-rat/repo"
	"psychic-rat/types"
)

////////////////////////////////////////////////////////////////////////////////
// CompanyApi implementations

////////////////////////////////////////////////////////////////////////////////
// Repo implementation

func NewCompanyApi(repos repo.Repos) *companyApiRepoImpl {
	return &companyApiRepoImpl{repos: repos}
}

type companyApiRepoImpl struct {
	repos repo.Repos
}

func (c *companyApiRepoImpl) GetCompanies() (types.CompanyListing, error) {
	companies := c.repos.Company.GetCompanies()
	results := types.CompanyListing{make([]types.CompanyElement, len(companies))}
	for i, co := range companies {
		results.Companies[i] = types.CompanyElement{co.Id, co.Name}
	}
	return results, nil
}

func (c *companyApiRepoImpl) GetById(id mdl.Id) (types.CompanyElement, error) {
	co, err := c.repos.Company.GetById(id)
	if err != nil {
		return types.CompanyElement{}, err
	}
	return types.CompanyElement{id, co.Name}, nil
}

////////////////////////////////////////////////////////////////////////////////
// RESTful implementation

func GetRestfulCompanyApi(baseUrl string) *restfulCompanyApi {
	return &restfulCompanyApi{baseUrl}
}

type restfulCompanyApi struct {
	url string
}

func (r *restfulCompanyApi) GetCompanies() (types.CompanyListing, error) {
	resp, err := http.Get(r.url + rest.CompanyApi)
	if err != nil {
		return types.CompanyListing{}, fmt.Errorf("get companies failed: %v", err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return types.CompanyListing{}, fmt.Errorf("error reading company response: %v", err)
	}
	return companiesFromJson(bytes)
}

func companiesFromJson(bytes []byte) (types.CompanyListing, error) {
	companies := types.CompanyListing{}
	if err := json.Unmarshal(bytes, &companies); err != nil {
		return companies, fmt.Errorf("failed to unmarshal companies: %v", err)
	}
	return companies, nil
}

func (r *restfulCompanyApi) GetById(id mdl.Id) (types.CompanyElement, error) {
	resp, err := http.Get(r.url + rest.CompanyApi + fmt.Sprintf("?company=%v", id))
	if err != nil {
		return types.CompanyElement{}, fmt.Errorf("could not retrieve company %v : %v", id, err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return types.CompanyElement{}, fmt.Errorf("error reading company element from response: %v", err)
	}
	return companyFromJson(bytes)
}

func companyFromJson(bytes []byte) (types.CompanyElement, error) {
	co := types.CompanyElement{}
	if err := json.Unmarshal(bytes, &co); err != nil {
		return co, fmt.Errorf("failed to unmarshal company: %v", err)
	}
	return co, nil
}

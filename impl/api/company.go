package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"psychic-rat/api/rest"
	"psychic-rat/mdl"
	"psychic-rat/types"
)

////////////////////////////////////////////////////////////////////////////////
// CompanyApi implementations

////////////////////////////////////////////////////////////////////////////////
// Repo implementation

func NewCompanyApi(repos Repos) *companyApiRepoImpl {
	return &companyApiRepoImpl{repos: repos}
}

type companyApiRepoImpl struct {
	repos Repos
}

func (c *companyApiRepoImpl) GetCompanies() (types.CompanyListing, error) {
	companies := c.repos.Company.GetCompanies()
	results := types.CompanyListing{make([]types.Company, len(companies))}
	for i, co := range companies {
		results.Companies[i] = types.Company{co.Id, co.Name}
	}
	return results, nil
}

func (c *companyApiRepoImpl) GetById(id mdl.ID) (types.Company, error) {
	co, err := c.repos.Company.GetById(id)
	if err != nil {
		return types.Company{}, err
	}
	return types.Company{id, co.Name}, nil
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

func (r *restfulCompanyApi) GetById(id mdl.ID) (types.Company, error) {
	resp, err := http.Get(r.url + rest.CompanyApi + fmt.Sprintf("?company=%v", id))
	if err != nil {
		return types.Company{}, fmt.Errorf("could not retrieve company %v : %v", id, err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return types.Company{}, fmt.Errorf("error reading company element from response: %v", err)
	}
	return companyFromJson(bytes)
}

func companyFromJson(bytes []byte) (types.Company, error) {
	co := types.Company{}
	if err := json.Unmarshal(bytes, &co); err != nil {
		return co, fmt.Errorf("failed to unmarshal company: %v", err)
	}
	return co, nil
}

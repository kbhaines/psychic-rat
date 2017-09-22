package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"psychic-rat/api/rest"
	"psychic-rat/mdl"
	"psychic-rat/repo/companyrepo"
)

////////////////////////////////////////////////////////////////////////////////
// CompanyApi implementations

////////////////////////////////////////////////////////////////////////////////
// Repo implementation

func getRepoCompanyApi() CompanyApi {
	return &companyApiRepoImpl{}
}

type companyApiRepoImpl struct{}

func (c *companyApiRepoImpl) GetCompanies() (CompanyListing, error) {
	companies := companyrepo.GetCompanyRepoMapImpl().GetCompanies()
	results := CompanyListing{make([]CompanyElement, len(companies))}
	for i, c := range companies {
		results.Companies[i] = CompanyElement{c.Id, c.Name}
	}
	return results, nil
}

func (c *companyApiRepoImpl) GetById(id mdl.Id) (CompanyElement, error) {
	co, err := companyrepo.GetCompanyRepoMapImpl().GetById(id)
	if err != nil {
		return CompanyElement{}, err
	}
	return CompanyElement{id, co.Name}, nil
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

func (r *restfulCompanyApi) GetById(id mdl.Id) (CompanyElement, error) {
	resp, err := http.Get(r.url + rest.CompanyApi + fmt.Sprintf("?company=%v", id))
	if err != nil {
		return CompanyElement{}, fmt.Errorf("could not retrieve company %v : %v", id, err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return CompanyElement{}, fmt.Errorf("error reading company element from response: %v", err)
	}
	return companyFromJson(bytes)
}

func companyFromJson(bytes []byte) (CompanyElement, error) {
	co := CompanyElement{}
	if err := json.Unmarshal(bytes, &co); err != nil {
		return co, fmt.Errorf("failed to unmarshal company: %v", err)
	}
	return co, nil
}

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

type Companies interface {
	Create(company mdl.CompanyRecord) error
	GetCompanies() []mdl.CompanyRecord
	GetById(mdl.Id) (*mdl.CompanyRecord, error)
}

type Items interface {
	Create(item mdl.ItemRecord) error
	GetById(id mdl.Id) (*mdl.ItemRecord, error)
	GetAllByCompany(companyId mdl.Id) []mdl.ItemRecord
	Update(id mdl.Id, item mdl.ItemRecord)
	List() []mdl.ItemRecord
}

type Pledges interface {
	Create(pledge mdl.PledgeRecord) error
	GetById(id mdl.Id) (*mdl.PledgeRecord, error)
	GetByUser(id mdl.Id) []mdl.Id
	List() []mdl.PledgeRecord
}

type Users interface {
	Create(user mdl.UserRecord) error
	GetById(id mdl.Id) (*mdl.UserRecord, error)
}

type Repos struct {
	Company Companies
	Item    Items
	Pledge  Pledges
	User    Users
}

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

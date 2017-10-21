package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"psychic-rat/api/rest"
	"psychic-rat/mdl"
	"psychic-rat/types"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

////////////////////////////////////////////////////////////////////////////////
// Implementations

////////////////////////////////////////////////////////////////////////////////
// Repo implementation

func NewPledgeApi(repos Repos) *pledgeApiRepoImpl {
	return &pledgeApiRepoImpl{repos: repos}
}

type pledgeApiRepoImpl struct {
	repos Repos
}

func (p *pledgeApiRepoImpl) NewPledge(itemId mdl.Id, userId mdl.Id) (newId mdl.Id, err error) {
	_, err = p.repos.Item.GetById(itemId)
	if err != nil {
		return newId, fmt.Errorf("error retrieving item %v: %v", itemId, err)
	}
	_, err = p.repos.User.GetById(userId)
	if userId != mdl.Id(0) && err != nil {
		return newId, fmt.Errorf("error retrieving user %v: %v", userId, err)
	}
	newPledge := mdl.PledgeRecord{Id: mdl.Id(uuid.NewV4().String()), UserId: userId, ItemId: itemId, Timestamp: time.Now()}
	p.repos.Pledge.Create(newPledge)
	return newPledge.Id, nil
}

func (p *pledgeApiRepoImpl) ListPledges() (types.PledgeListing, error) {
	return types.PledgeListing{}, nil
}

////////////////////////////////////////////////////////////////////////////////
// Restful implementation

func GetRestfulPledgeApiImpl(url string) *restPledgeApiImpl {
	return &restPledgeApiImpl{url}
}

type restPledgeApiImpl struct {
	url string
}

type pledgeResponse struct {
	Id mdl.Id `json:"id"`
}

func NewPledgeRequest(itemId mdl.Id) types.PledgeRequest {
	return types.PledgeRequest{itemId}
}

func (a *restPledgeApiImpl) NewPledge(itemId mdl.Id, userId mdl.Id) (mdl.Id, error) {
	jsonString := rest.ToJsonString(NewPledgeRequest(itemId))
	body := strings.NewReader(jsonString)
	resp, err := http.Post(a.url+rest.PledgeApi, "application/json", body)
	if resp.StatusCode != http.StatusOK {
		return mdl.Id(0), fmt.Errorf("request returned: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	r := pledgeResponse{}
	err = json.Unmarshal(bytes, &r)
	if err != nil {
		return mdl.Id(0), fmt.Errorf("unable to decode response: %v", err)
	}
	return r.Id, nil
}

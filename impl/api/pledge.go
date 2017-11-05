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

func (p *pledgeApiRepoImpl) NewPledge(itemId mdl.ID, userId mdl.ID) (newId mdl.ID, err error) {
	_, err = p.repos.Item.GetById(itemId)
	if err != nil {
		return newId, fmt.Errorf("error retrieving item %v: %v", itemId, err)
	}
	_, err = p.repos.User.GetById(userId)
	if userId != mdl.ID(0) && err != nil {
		return newId, fmt.Errorf("error retrieving user %v: %v", userId, err)
	}
	newPledge := mdl.Pledge{Id: mdl.ID(uuid.NewV4().String()), UserID: userId, ItemID: itemId, Timestamp: time.Now()}
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
	Id mdl.ID `json:"id"`
}

func NewPledgeRequest(itemId mdl.ID) types.PledgeRequest {
	return types.PledgeRequest{itemId}
}

func (a *restPledgeApiImpl) NewPledge(itemId mdl.ID, userId mdl.ID) (mdl.ID, error) {
	jsonString := rest.ToJsonString(NewPledgeRequest(itemId))
	body := strings.NewReader(jsonString)
	resp, err := http.Post(a.url+rest.PledgeApi, "application/json", body)
	if resp.StatusCode != http.StatusOK {
		return mdl.ID(0), fmt.Errorf("request returned: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	r := pledgeResponse{}
	err = json.Unmarshal(bytes, &r)
	if err != nil {
		return mdl.ID(0), fmt.Errorf("unable to decode response: %v", err)
	}
	return r.Id, nil
}

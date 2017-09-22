package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"psychic-rat/api/rest"
	"psychic-rat/mdl"
	"psychic-rat/repo/itemrepo"
	"psychic-rat/repo/pledgerepo"
	"psychic-rat/repo/userrepo"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

type PledgeApi interface {
	NewPledge(itemId mdl.Id, userId mdl.Id) (mdl.Id, error)
	//ListPledges() PledgeListing
}

type PledgeListing struct {
	Pledges []PledgeElement `json:"pledges"`
}

type PledgeElement struct {
	PledgeId  mdl.Id      `json:"id"`
	Item      ItemElement `json:"item"`
	Timestamp time.Time   `json:"timestamp"`
}

type PledgeRequest struct {
	ItemId mdl.Id `json:"itemId"`
}

////////////////////////////////////////////////////////////////////////////////
// Implementations

////////////////////////////////////////////////////////////////////////////////
// Repo implementation

type pledgeApiRepoImpl struct{}

func GetRepoPledgeApiImpl() PledgeApi {
	return &pledgeApiRepoImpl{}
}

func (p *pledgeApiRepoImpl) NewPledge(itemId mdl.Id, userId mdl.Id) (newId mdl.Id, err error) {
	_, err = itemrepo.GetItemRepoMapImpl().GetById(itemId)
	if err != nil {
		return newId, fmt.Errorf("error retrieving item %v: %v", itemId, err)
	}
	_, err = userrepo.GetUserRepoMapImpl().GetById(userId)
	if userId != mdl.Id(0) && err != nil {
		return newId, fmt.Errorf("error retrieving user %v: %v", userId, err)
	}
	newPledge := mdl.PledgeRecord{Id: mdl.Id(uuid.NewV4().String()), UserId: userId, ItemId: itemId, Timestamp: time.Now()}
	pledgerepo.GetPledgeRepoMapImpl().Create(newPledge)
	return newPledge.Id, nil
}

func (p *pledgeApiRepoImpl) ListPledges() (PledgeListing, error) {
	return PledgeListing{}, nil
}

////////////////////////////////////////////////////////////////////////////////
// Restful implementation

func GetRestfulPledgeApiImpl(url string) PledgeApi {
	return &restPledgeApiImpl{url}
}

type restPledgeApiImpl struct {
	url string
}

type pledgeResponse struct {
	Id mdl.Id `json:"id"`
}

func NewPledgeRequest(itemId mdl.Id) PledgeRequest {
	return PledgeRequest{itemId}
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

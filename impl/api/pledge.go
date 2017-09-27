package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	a "psychic-rat/api"
	"psychic-rat/api/rest"
	"psychic-rat/impl/repo"
	"psychic-rat/mdl"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

////////////////////////////////////////////////////////////////////////////////
// Implementations

////////////////////////////////////////////////////////////////////////////////
// Repo implementation

type pledgeApiRepoImpl struct{}

func getRepoPledgeApi() a.PledgeApi {
	return &pledgeApiRepoImpl{}
}

func (p *pledgeApiRepoImpl) NewPledge(itemId mdl.Id, userId mdl.Id) (newId mdl.Id, err error) {
	_, err = repo.Get().Pledge.GetById(itemId)
	if err != nil {
		return newId, fmt.Errorf("error retrieving item %v: %v", itemId, err)
	}
	_, err = repo.Get().User.GetById(userId)
	if userId != mdl.Id(0) && err != nil {
		return newId, fmt.Errorf("error retrieving user %v: %v", userId, err)
	}
	newPledge := mdl.PledgeRecord{Id: mdl.Id(uuid.NewV4().String()), UserId: userId, ItemId: itemId, Timestamp: time.Now()}
	repo.Get().Pledge.Create(newPledge)
	return newPledge.Id, nil
}

func (p *pledgeApiRepoImpl) ListPledges() (a.PledgeListing, error) {
	return a.PledgeListing{}, nil
}

////////////////////////////////////////////////////////////////////////////////
// Restful implementation

func GetRestfulPledgeApiImpl(url string) a.PledgeApi {
	return &restPledgeApiImpl{url}
}

type restPledgeApiImpl struct {
	url string
}

type pledgeResponse struct {
	Id mdl.Id `json:"id"`
}

func NewPledgeRequest(itemId mdl.Id) a.PledgeRequest {
	return a.PledgeRequest{itemId}
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

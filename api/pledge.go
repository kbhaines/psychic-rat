package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"psychic-rat/ctr"
	"psychic-rat/mdl/item"
	"psychic-rat/mdl/pledge"
	"psychic-rat/mdl/user"
	"time"
)

import (
	"log"
	"sort"
)

type PledgeListing struct {
	Pledges []PledgeElement `json:"pledges"`
}

type PledgeElement struct {
	PledgeId  pledge.Id   `json:"id"`
	Item      ItemElement `json:"item"`
	Timestamp time.Time   `json:"timestamp"`
}

type PledgeRequest struct {
	ItemId item.Id `json:"itemId"`
}

func PledgeHandler(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		handlePost(writer, request)

	case http.MethodGet:
		handleGet(writer, request)

	default:
		unsupportedMethod(writer)
	}
}

func NewPledgeRequest(itemId item.Id) PledgeRequest {
	return PledgeRequest{itemId}
}

func handlePost(writer http.ResponseWriter, request *http.Request) {
	pledge := PledgeRequest{}

	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		logInternalError(writer, err)
		return
	}
	err = json.Unmarshal(body, &pledge)
	if err != nil {
		logInternalError(writer, err)
		return
	}
	userId := getCurrentUserId()
	_, err = ctr.GetController().Pledge().AddPledge(pledge.ItemId, userId)
	if err != nil {
		logInternalError(writer, err)
		return
	}
	writeUserPledges(writer, userId)
}

func writeUserPledges(writer http.ResponseWriter, userId user.Id) {
	userPledges := getUserPledges(userId)
	log.Printf("pledges: %v", userPledges.Pledges)
	json, err := json.Marshal(userPledges)
	if err != nil {
		logInternalError(writer, err)
		return
	}
	fmt.Fprintf(writer, "%s", json)
}

func (p *PledgeElement) String() string { return fmt.Sprintf("id:%v time:%v", p.PledgeId, p.Timestamp) }

func returnIfElse(b bool, ifTrue, ifFalse interface{}) interface{} {
	if b {
		return ifTrue
	} else {
		return ifFalse
	}
}

func getUserPledges(id user.Id) PledgeListing {
	ps := ctr.GetController().Pledge().ListPledges(func(p pledge.Record) pledge.Record {
		return ifElse(id == p.UserId(), p, nil).(pledge.Record)
	})
	sort.Sort(pledge.ByTimeStamp(ps))
	ps2 := make([]PledgeElement, len(ps))
	for i, p := range ps {
		item, err := ctr.GetController().Item().GetById(p.ItemId())
		if err != nil {
			panic(fmt.Sprintf("data inconsistency error %v. item %v for pledge %v does not exist", err, p.ItemId(), p.Id()))
		}
		ps2[i] = PledgeElement{p.Id(), ItemElement{item.Id(), item.Make(), item.Model()}, p.TimeStamp()}
	}
	return PledgeListing{ps2}
}

func handleGet(writer http.ResponseWriter, request *http.Request) {
	userId := getCurrentUserId()
	writeUserPledges(writer, userId)
}

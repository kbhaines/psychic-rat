package api

import (
	"net/http"
	"fmt"
	"psychic-rat/ctr"
	"psychic-rat/mdl/item"
	"psychic-rat/mdl/user"
	"encoding/json"
	"io/ioutil"
	"psychic-rat/mdl/pledge"
	"time"
)

import (
	"log"
	"sort"
)

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

type pledgeRequest struct {
	ItemId item.Id `json:"itemId"`
}

func handlePost(writer http.ResponseWriter, request *http.Request) {
	pledge := pledgeRequest{}

	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		logInternalError(writer, err);
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

type pledgeListing struct {
	Pledges []pledge.Record `json:"pledges"`
}

type pledgeElement struct {
	PledgeId  pledge.Id `json:"id"`
	Item      item.Record `json:"item"`
	Timestamp time.Time `json:"timestamp"`
}

func (p *pledgeElement) Id() pledge.Id        { return p.Id() }
func (p *pledgeElement) UserId() user.Id      { return p.UserId() }
func (p *pledgeElement) ItemId() item.Id      { return item.Id(0) }
func (p *pledgeElement) TimeStamp() time.Time { return p.Timestamp }

func (p *pledgeElement) String() string { return fmt.Sprintf("id:%v time:%v", p.PledgeId, p.Timestamp) }

func returnIfElse(b bool, ifTrue, ifFalse interface{}) interface{} {
	if b {
		return ifTrue
	} else {
		return ifFalse
	}
}

func getUserPledges(id user.Id) pledgeListing {
	ps := ctr.GetController().Pledge().ListPledges(func(p pledge.Record) pledge.Record {
		if id == p.UserId() {
			item, err := ctr.GetController().Item().GetById(p.ItemId())
			if err != nil {
				panic(err)
			}
			i := &itemElement{item.Id(), item.Make(), item.Model()}
			return &pledgeElement{p.Id(), i, p.TimeStamp()}
		} else {
			return nil
		}
	})
	sort.Sort(pledge.ByTimeStamp(ps))
	return pledgeListing{ps}
}

func handleGet(writer http.ResponseWriter, request *http.Request) {
	userId := getCurrentUserId()
	writeUserPledges(writer, userId)
}

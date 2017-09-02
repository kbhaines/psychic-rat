package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"psychic-rat/api/client/ctr"
	"psychic-rat/mdl/company"
	"psychic-rat/mdl/item"
)

func CompanyHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		unsupportedMethod(writer)
		return
	}
	ToJson(writer, ctr.GetCompanies())
}

//// ITEMS

func ItemHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		unsupportedMethod(writer)
		return
	}
	if err := request.ParseForm(); err != nil {
		unableToParseForm(err, writer)
		return
	}

	companyId := request.Form.Get("company")
	if companyId == "" {
		errorResponse(writer, fmt.Errorf("no company specified"))
		return
	}
	json, err := json.Marshal(getItemsForCompany(company.Id(companyId)))
	if err != nil {
		logInternalError(writer, err)
		return
	}
	fmt.Fprintf(writer, "%s", json)
}

func ifElse(b bool, t, f interface{}) interface{} {
	if b {
		return t
	}
	return f
}

func getItemsForCompany(companyId company.Id) ItemReport {
	is := ctr.GetController().Item().ListItems(func(i item.Record) item.Record {
		if companyId == i.Company() {
			return i
		} else {
			return nil
		}
	})
	report := ItemReport{make([]ItemElement, len(is))}
	for i, v := range is {
		report.Items[i] = ItemElement{v.Id(), v.Make(), v.Model()}
	}
	return report
}

func ItemsFromJson(bytes []byte) (ItemReport, error) {
	var items ItemReport
	if err := json.Unmarshal(bytes, &items); err != nil {
		return items, fmt.Errorf("failed to unmarshal items: %v", err)
	}
	return items, nil
}

// PLEDGES

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

package main

import (
	"net/http"
	"net/http/httptest"
	"psychic-rat/api"
	"psychic-rat/mdl"
	"reflect"
	"testing"
)

var mockItemReport = api.ItemReport{
	Items: []api.ItemElement{
		api.ItemElement{Id: mdl.Id("123"), Make: "phone", Model: "x124", Company: "bigco1"},
		api.ItemElement{Id: mdl.Id("124"), Make: "phone", Model: "x125", Company: "bigco2"},
		api.ItemElement{Id: mdl.Id("125"), Make: "phone", Model: "x126", Company: "bigco2"},
		api.ItemElement{Id: mdl.Id("126"), Make: "phone", Model: "x127", Company: "bigco3"},
	},
}

type mockItemApi struct{}

func (m *mockItemApi) ListItems() (api.ItemReport, error) {
	return mockItemReport, nil
}

func (m *mockItemApi) GetById(id mdl.Id) (api.ItemElement, error) {
	panic("not implemented")
}

func mockSession(request *http.Request) bool {
	return true
}

func TestPledgeListItems(t *testing.T) {
	apis = api.Api{Item: &mockItemApi{}}
	expectedVars := pageVariables{Username: "Kevin", Items: mockItemReport.Items}
	renderPage = getRenderMock(t, "pledge.html.tmpl", expectedVars)
	isUserLoggedIn = mockSession
	req := &http.Request{Method: "GET"}
	PledgePageHandler(nil, req)
}

func getRenderMock(t *testing.T, expectedTemplate string, expectedVars pageVariables) renderFunc {
	return func(writer http.ResponseWriter, templateName string, templateVars interface{}) {
		t.Logf("Rendering %s with %v", templateName, templateVars)
		if templateName != expectedTemplate {
			t.Errorf("wrong template, got %s but wanted %s", templateName, expectedTemplate)
		}
		if v, ok := templateVars.(pageVariables); !ok {
			t.Errorf("did not match type")
		} else if !reflect.DeepEqual(v, expectedVars) {
			t.Errorf("wrong variables, got %v but wanted %v", v, expectedVars)
		}
	}
}

func TestHomePage(t *testing.T) {
	apis = api.Api{Item: &mockItemApi{}}
	expectedVars := pageVariables{Username: "Kevin"}
	renderPage = getRenderMock(t, "home.html.tmpl", expectedVars)
	isUserLoggedIn = mockSession
	req := &http.Request{Method: "GET"}
	HomePageHandler(nil, req)
}

var protectedPages = []handlerFunc{
	PledgePageHandler,
	ThanksPageHandler,
}

func TestProtectedPages(t *testing.T) {
	req := &http.Request{Method: "GET"}
	isUserLoggedIn = func(r *http.Request) bool { return false }
	for _, page := range protectedPages {
		writer := httptest.NewRecorder()
		page(writer, req)
		result := writer.Result()
		if result.StatusCode != http.StatusForbidden {
			t.Errorf("testing %v, got %v, but wanted %v", page, result.StatusCode, http.StatusForbidden)
		}
	}
}

func TestPledgeAuthFailure(t *testing.T) {
	isUserLoggedIn = func(r *http.Request) bool { return false }

	for _, method := range []string{"GET", "POST"} {
		writer := httptest.NewRecorder()
		req := &http.Request{Method: method}
		PledgePageHandler(writer, req)
		result := writer.Result()
		if result.StatusCode != http.StatusForbidden {
			t.Errorf("testing %v, got %v, but wanted %v", method, result.StatusCode, http.StatusForbidden)
		}
	}
}

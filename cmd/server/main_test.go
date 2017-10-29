// +build !integration

package main

import (
	"net/http"
	"net/http/httptest"
	"psychic-rat/mdl"
	"psychic-rat/types"
	"reflect"
	"testing"
)

var mockItemReport = types.ItemReport{
	Items: []types.ItemElement{
		types.ItemElement{Id: mdl.Id("123"), Make: "phone", Model: "x124", Company: "bigco1"},
		types.ItemElement{Id: mdl.Id("124"), Make: "phone", Model: "x125", Company: "bigco2"},
		types.ItemElement{Id: mdl.Id("125"), Make: "phone", Model: "x126", Company: "bigco2"},
		types.ItemElement{Id: mdl.Id("126"), Make: "phone", Model: "x127", Company: "bigco3"},
	},
}

type mockItemApi struct{}

func (m *mockItemApi) ListItems() (types.ItemReport, error) {
	return mockItemReport, nil
}

func (m *mockItemApi) GetById(id mdl.Id) (types.ItemElement, error) {
	panic("not implemented")
}

func mockSession(request *http.Request) bool {
	return true
}

func TestPledgeListItems(t *testing.T) {
	apis = Api{Item: &mockItemApi{}}
	expectedVars := pageVariables{Items: mockItemReport.Items}
	renderPage = getRenderMock(t, "pledge.html.tmpl", expectedVars)
	isUserLoggedIn = mockSession
	req := &http.Request{Method: "GET"}
	PledgePageHandler(nil, req)
}

func getRenderMock(t *testing.T, expectedTemplate string, expectedVars pageVariables) renderFunc {
	return func(writer http.ResponseWriter, templateName string, templateVars *pageVariables) {
		t.Logf("Rendering %s with %v", templateName, templateVars)
		if templateName != expectedTemplate {
			t.Errorf("wrong template, got %s but wanted %s", templateName, expectedTemplate)
		}
		if !reflect.DeepEqual(*templateVars, expectedVars) {
			t.Errorf("wrong variables, wanted %v, got %v", expectedVars, *templateVars)
		}
	}
}

func TestHomePage(t *testing.T) {
	apis = Api{}
	expectedVars := pageVariables{}
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

package web

import (
	"net/http"
	"net/http/httptest"
	"psychic-rat/types"
	"reflect"
	"testing"
)

var (
	mockCompanies = []types.Company{
		types.Company{Name: "bigco1"},
		types.Company{Name: "bigco2"},
		types.Company{Name: "bigco3"},
	}

	mockItemReport = []types.Item{
		types.Item{Id: 123, Make: "phone", Model: "x124", Company: mockCompanies[0]},
		types.Item{Id: 124, Make: "phone", Model: "x125", Company: mockCompanies[1]},
		types.Item{Id: 125, Make: "phone", Model: "x126", Company: mockCompanies[1]},
		types.Item{Id: 126, Make: "phone", Model: "x127", Company: mockCompanies[2]},
	}
)

type mockItemApi struct{}

func (m *mockItemApi) ListItems() ([]types.Item, error) {
	return mockItemReport, nil
}

func (m *mockItemApi) GetItem(id int) (types.Item, error) {
	panic("not implemented")
}

func mockSession(request *http.Request) bool {
	return true
}

func TestPledgeListItems(t *testing.T) {
	apis = API{Item: &mockItemApi{}}
	expectedVars := pageVariables{Items: mockItemReport}
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

func TestHomePageTemplate(t *testing.T) {
	apis = API{}
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
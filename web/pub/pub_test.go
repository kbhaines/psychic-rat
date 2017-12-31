package pub

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
		types.Item{ID: 123, Make: "phone", Model: "x124", Company: mockCompanies[0]},
		types.Item{ID: 124, Make: "phone", Model: "x125", Company: mockCompanies[1]},
		types.Item{ID: 125, Make: "phone", Model: "x126", Company: mockCompanies[1]},
		types.Item{ID: 126, Make: "phone", Model: "x127", Company: mockCompanies[2]},
	}

	protectedPages = []http.HandlerFunc{
		PledgePageHandler,
		ThanksPageHandler,
	}
)

type mockItemAPI struct{}

type mockAuthHandler struct {
	user *types.User
}

type mockRenderer struct {
	expectedTemplate string
	expectedVars     pageVariables
	t                *testing.T
}

func TestPledgeListItems(t *testing.T) {
	apis = APIS{Item: &mockItemAPI{}}
	authHandler = &mockAuthHandler{user: &types.User{}}

	expectedVars := pageVariables{Items: mockItemReport}
	renderer = getRenderMock(t, "pledge.html.tmpl", expectedVars)
	req := &http.Request{Method: "GET"}
	PledgePageHandler(nil, req)
}

func getRenderMock(t *testing.T, expectedTemplate string, expectedVars pageVariables) Renderer {
	return &mockRenderer{
		expectedTemplate: expectedTemplate,
		expectedVars:     expectedVars,
		t:                t,
	}
}

func TestHomePageTemplate(t *testing.T) {
	apis = APIS{}
	expectedVars := pageVariables{}
	renderer = getRenderMock(t, "home.html.tmpl", expectedVars)
	req := &http.Request{Method: "GET"}
	HomePageHandler(nil, req)
}

func TestProtectedPages(t *testing.T) {
	renderer = getRenderMock(t, "", pageVariables{})
	apis = APIS{Item: &mockItemAPI{}}
	authHandler = &mockAuthHandler{}
	req := &http.Request{Method: "GET"}
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
func (m *mockItemAPI) ListItems() ([]types.Item, error)   { return mockItemReport, nil }
func (m *mockItemAPI) GetItem(id int) (types.Item, error) { panic("not implemented") }

func (a *mockAuthHandler) Handler(http.ResponseWriter, *http.Request)         {}
func (a *mockAuthHandler) GetLoggedInUser(*http.Request) (*types.User, error) { return a.user, nil }

func (r *mockRenderer) Render(w http.ResponseWriter, templateName string, vars interface{}) error {
	t := r.t
	t.Logf("Rendering %s with %v", templateName, vars)
	if templateName != r.expectedTemplate {
		t.Errorf("wrong template, got %s but wanted %s", templateName, r.expectedTemplate)
	}
	if !reflect.DeepEqual(*vars.(*pageVariables), r.expectedVars) {
		t.Errorf("wrong variables, wanted %v, got %v", r.expectedVars, *vars.(*pageVariables))
	}
	return nil
}

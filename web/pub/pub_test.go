package pub

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
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

	mockItemList = []types.Item{
		types.Item{ID: 123, Make: "phone", Model: "x124", Company: mockCompanies[0], USDValue: 100},
		types.Item{ID: 124, Make: "phone", Model: "x125", Company: mockCompanies[1], USDValue: 100},
		types.Item{ID: 125, Make: "phone", Model: "x126", Company: mockCompanies[1], USDValue: 100},
		types.Item{ID: 126, Make: "phone", Model: "x127", Company: mockCompanies[2], USDValue: 100},
	}

	protectedPages = []http.HandlerFunc{
		PledgePageHandler,
		ThanksPageHandler,
	}
)

type mockItemAPI struct{}
type mockPledgeAPI struct {
	userID string
	itemID int
	value  int
	t      *testing.T
}

type mockNewItemsAPI struct {
	newItem *types.NewItem
	t       *testing.T
}

type mockAuthHandler struct {
	user *types.User
}

type mockRenderer struct {
	expectedTemplate string
	expectedVars     pageVariables
	t                *testing.T
}

func TestPledgeListItems(t *testing.T) {
	itemsAPI = &mockItemAPI{}
	authHandler = &mockAuthHandler{user: &types.User{}}

	expectedVars := pageVariables{Items: mockItemList}
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
	expectedVars := pageVariables{}
	renderer = getRenderMock(t, "home.html.tmpl", expectedVars)
	req := &http.Request{Method: "GET"}
	HomePageHandler(httptest.NewRecorder(), req)
}

func TestProtectedPages(t *testing.T) {
	renderer = getRenderMock(t, "", pageVariables{})
	itemsAPI = &mockItemAPI{}
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

func TestPledgePost(t *testing.T) {
	values := url.Values{"item": {"1"}}
	itemsAPI = &mockItemAPI{}
	pledgeAPI = &mockPledgeAPI{userID: "test1", itemID: 1, value: 100, t: t}
	expectedVars := pageVariables{User: types.User{ID: "test1"}, Items: []types.Item{mockItemList[1]}}
	renderer = getRenderMock(t, "pledge-post.html.tmpl", expectedVars)
	authHandler = &mockAuthHandler{user: &types.User{ID: "test1"}}
	req := &http.Request{Method: "POST", PostForm: values, URL: &url.URL{Path: "/pledge"}}
	writer := httptest.NewRecorder()
	PledgePageHandler(writer, req)
	if writer.Result().StatusCode != http.StatusSeeOther {
		t.Fatalf("expected 303, got %s", writer.Result().Status)
	}
	pledgeAPI.(*mockPledgeAPI).checkUsed()
}

func TestNewItemPledgePost(t *testing.T) {
	newItem := struct {
		company  string
		make     string
		model    string
		currency string
		value    string
	}{"newco", "newmake", "newmodel", "1", "100"}

	values := url.Values{"company": {newItem.company}, "make": {newItem.make}, "model": {newItem.model}, "currencyID": {newItem.currency}, "value": {newItem.value}}
	itemsAPI = &mockItemAPI{}
	newItemsAPI = &mockNewItemsAPI{
		t: t,
		newItem: &types.NewItem{
			Company:    newItem.company,
			Make:       newItem.make,
			Model:      newItem.model,
			UserID:     "test1",
			CurrencyID: 1,
			Value:      100,
			IsPledge:   true,
		},
	}
	expectedVars := pageVariables{
		User: types.User{ID: "test1"},
		Items: []types.Item{
			types.Item{Make: newItem.make, Model: newItem.model, Company: types.Company{Name: newItem.company}},
		},
	}
	renderer = getRenderMock(t, "pledge-post-new-item.html.tmpl", expectedVars)
	authHandler = &mockAuthHandler{user: &types.User{ID: "test1"}}
	req := &http.Request{Method: "POST", PostForm: values}
	writer := httptest.NewRecorder()
	NewItemHandler(writer, req)
	if writer.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected OK, got %s", writer.Result().Status)
	}
	newItemsAPI.(*mockNewItemsAPI).checkUsed()
	renderer.(*mockRenderer).checkUsed()
}

func (m *mockItemAPI) ListItems() ([]types.Item, error)          { return mockItemList, nil }
func (m *mockItemAPI) GetItem(id int) (types.Item, error)        { return mockItemList[id], nil }
func (m *mockItemAPI) ListCurrencies() ([]types.Currency, error) { return nil, nil }

func (p *mockPledgeAPI) AddPledge(itemID int, userID string, value int) (*types.Pledge, error) {
	p.t.Helper()
	if p.userID != userID || p.itemID != itemID || p.value != value {
		err := fmt.Errorf("expected AddPledge to be called with userID:%s itemID:%d value:%d, got %s %d %d", p.userID, p.itemID, p.value, userID, itemID, value)
		p.t.Fatal(err)
		return nil, err
	}
	p.userID = ""
	return &types.Pledge{Item: mockItemList[itemID], PledgeID: 1, UserID: userID}, nil
}

func (p *mockPledgeAPI) checkUsed() {
	p.t.Helper()
	if p.userID != "" {
		p.t.Fatalf("expected call to AddPledge with %d,%s, didn't happen", p.itemID, p.userID)
	}
}

func (a *mockAuthHandler) Handler(http.ResponseWriter, *http.Request)         {}
func (a *mockAuthHandler) GetLoggedInUser(*http.Request) (*types.User, error) { return a.user, nil }
func (a *mockAuthHandler) LogOut(http.ResponseWriter, *http.Request) error    { return nil }

func (r *mockRenderer) Render(w http.ResponseWriter, templateName string, vars interface{}) error {
	t := r.t
	t.Helper()
	t.Logf("Rendering %s with %v", templateName, vars)
	if templateName != r.expectedTemplate {
		t.Errorf("wrong template, got %s but wanted %s", templateName, r.expectedTemplate)
	}
	if !reflect.DeepEqual(*vars.(*pageVariables), r.expectedVars) {
		t.Errorf("wrong variables, wanted %v, got %v", r.expectedVars, *vars.(*pageVariables))
	}
	r.expectedTemplate = ""
	return nil
}

func (r *mockRenderer) checkUsed() {
	r.t.Helper()
	if r.expectedTemplate != "" {
		r.t.Fatalf("expected call to Render with template named %s and vars %v didn't occur", r.expectedTemplate, r.expectedVars)
	}
}

func (n *mockNewItemsAPI) AddNewItem(item types.NewItem) (*types.NewItem, error) {
	n.t.Helper()
	if n.newItem == nil {
		err := fmt.Errorf("unexpected call to AddNewItem with params %v", item)
		n.t.Fatal(err)
		return nil, err
	}
	if *n.newItem != item {
		err := fmt.Errorf("expected AddNewItem to call with %v, got %v", n.newItem, item)
		n.t.Fatal(err)
		return nil, err
	}
	n.newItem = nil
	return &item, nil
}

func (n *mockNewItemsAPI) checkUsed() {
	n.t.Helper()
	if n.newItem != nil {
		n.t.Fatalf("expected call to AddNewItem with %v, didn't happen", n.newItem)
	}
}

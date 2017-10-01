package main

import (
	"net/http"
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

func TestPledgeListItems(t *testing.T) {
	apis = api.Api{Item: &mockItemApi{}}
	expectedVars := variables{Username: "Kevin", Items: mockItemReport.Items}
	renderPage = getRenderFunc(t, "pledge.html.tmpl", expectedVars)
	PledgePageHandler(nil, nil)
}

func getRenderFunc(t *testing.T, expectedTemplate string, expectedVars variables) renderFunction {
	return func(writer http.ResponseWriter, templateName string, templateVars interface{}) {
		t.Logf("Rendering %s with %v", templateName, templateVars)
		if templateName != expectedTemplate {
			t.Errorf("wrong template, got %s but wanted %s", templateName, expectedTemplate)
		}
		if v, ok := templateVars.(variables); !ok {
			t.Errorf("did not match type")
		} else if !reflect.DeepEqual(v, expectedVars) {
			t.Errorf("wrong variables, got %v but wanted %v", v, expectedVars)
		}
	}
}

func TestHomePage(t *testing.T) {
	apis = api.Api{Item: &mockItemApi{}}
	expectedVars := variables{Username: "Kevin"}
	renderPage = getRenderFunc(t, "home.html.tmpl", expectedVars)
	HomePageHandler(nil, nil)
}

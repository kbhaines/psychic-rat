package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestHomePage(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	resp, err := http.Get(testUrl + "/")
	testPageStatus(resp, err, http.StatusOK, t)
}

func TestPledgeWithoutLogin(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	resp, err := http.Get(testUrl + "/pledge")
	testPageStatus(resp, err, http.StatusForbidden, t)
}

func TestThankYouWithoutLogin(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	resp, err := http.Get(testUrl + "/thanks")
	testPageStatus(resp, err, http.StatusForbidden, t)
}

func TestSignin(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	loginUser("test1", t)
}

func TestPledgeWithLogin(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)

	cookie := loginUser("test1", t)
	client := http.Client{Jar: cookie}
	resp, err := client.Get(testUrl + "/pledge")
	testPageStatus(resp, err, http.StatusOK, t)

	strBody := readResponseBody(resp, t)
	expected := []string{
		"user1 full",
		"<select ",
		"<input type=\"submit\"",
	}
	testStrings(strBody, expected, t)
}

func TestHappyPathPledge(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	cookie := loginUser("test1", t)
	client := http.Client{Jar: cookie}
	csrf := getCSRFToken(client, testUrl+"/pledge", t)
	data := url.Values{"item": {"1"}, "csrf": {csrf}}
	resp, err := client.PostForm(testUrl+"/pledge", data)
	testPageStatus(resp, err, http.StatusOK, t)
	if actual := resp.Request.URL.String(); actual != testUrl+"/thanks" {
		t.Fatalf("expected to land at /thanks, got %s", actual)
	}
}

func getCSRFToken(client http.Client, url string, t *testing.T) string {
	t.Helper()
	resp, err := client.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		t.Fatal(err)
	}
	token := doc.Find("input[name=csrf]").AttrOr("value", "")
	return token
}

func TestNewItem(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	cookie := loginUser("test1", t)
	client := http.Client{Jar: cookie}

	values := url.Values{"company": {"newCo1"}, "make": {"newmake"}, "model": {"newmodel"}, "currencyID": {"1"}, "value": {"100"}}
	values.Add("csrf", getCSRFToken(client, testUrl+"/pledge", t))
	resp, err := client.PostForm(testUrl+"/newitem", values)
	testPageStatus(resp, err, http.StatusOK, t)
}

func TestBadNewItems(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	cookie := loginUser("test1", t)
	client := http.Client{Jar: cookie}

	values := []url.Values{
		url.Values{"company": {""}, "make": {"newmake"}, "model": {"newmodel"}, "currencyID": {"1"}, "value": {"100"}},
		url.Values{"company": {"newCo1"}, "make": {""}, "model": {"newmodel"}, "currencyID": {"1"}, "value": {"100"}},
		url.Values{"company": {"newCo1"}, "make": {"newmake"}, "model": {""}, "currencyID": {"1"}, "value": {"100"}},
		url.Values{"company": {"newCo1"}, "make": {"newmake"}, "model": {"newmodel"}, "currencyID": {""}, "value": {"100"}},
		url.Values{"company": {"newCo1"}, "make": {"newmake"}, "model": {"newmodel"}, "currencyID": {"1"}, "value": {""}},
		url.Values{"company": {"newCo1"}, "make": {"newmake"}, "model": {"newmodel"}, "currencyID": {"xxx"}, "value": {"100"}},
		url.Values{"company": {"newCo1"}, "make": {"newmake"}, "model": {"newmodel"}, "currencyID": {"1"}, "value": {"xxx"}},
		url.Values{"company": {"newCo1"}, "make": {"newmake"}, "model": {"newmodel"}, "currencyID": {"999"}, "value": {"100"}},
	}

	for _, d := range values {
		d.Add("csrf", getCSRFToken(client, testUrl+"/pledge", t))
		resp, err := client.PostForm(testUrl+"/newitem", d)
		testPageStatus(resp, err, http.StatusBadRequest, t)
	}
	items, err := db.ListNewItems()
	if err != nil {
		t.Fatal(err)
	}
	if 0 != len(items) {
		t.Fatalf("expected 0 new items, got %d", len(items))
	}
}

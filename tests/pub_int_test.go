package tests

import (
	"net/http"
	"net/url"
	"testing"
)

func TestHomePage(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	resp, err := http.Get(testURL + "/")
	testPageStatus(resp, err, http.StatusOK, t)
}

func TestPledgeWithoutLogin(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	resp, err := http.Get(testURL + "/pledge")
	testPageStatus(resp, err, http.StatusOK, t)
	if resp.Request.URL.RequestURI() != "/signin" {
		t.Fatal("expected redirect to /signin")
	}
}

func TestThankYouWithoutLogin(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	resp, err := http.Get(testURL + "/thanks")
	testPageStatus(resp, err, http.StatusOK, t)
	if resp.Request.URL.RequestURI() != "/signin" {
		t.Fatal("expected redirect to /signin")
	}
}

func TestPledgeWithLogin(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)

	resp, err := getAuthdClient("user1", t).Get(testURL + "/pledge")
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
	client := getAuthdClient("admin", t)
	csrf := getCSRFToken(client, testURL+"/pledge", t)
	data := url.Values{"item": {"1"}, "csrf": {csrf}}
	resp, err := client.PostForm(testURL+"/pledge", data)
	testPageStatus(resp, err, http.StatusOK, t)
	if actual := resp.Request.URL.String(); actual != testURL+"/thanks" {
		t.Fatalf("expected to land at /thanks, got %s", actual)
	}
}

func TestNewItem(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	values := url.Values{"company": {"newCo1"}, "make": {"newmake"}, "model": {"newmodel"}, "currencyID": {"1"}, "value": {"100"}}
	client := getAuthdClient("admin", t)
	values.Add("csrf", getCSRFToken(client, testURL+"/pledge", t))
	resp, err := client.PostForm(testURL+"/newitem", values)
	testPageStatus(resp, err, http.StatusOK, t)
}

func TestBadNewItems(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)

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

	client := getAuthdClient("admin", t)
	for _, d := range values {
		d.Add("csrf", getCSRFToken(client, testURL+"/pledge", t))
		resp, err := client.PostForm(testURL+"/newitem", d)
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

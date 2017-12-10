package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var testUrl string

func TestHomePage(t *testing.T) {
	server := newServer(t)
	defer server.Close()
	resp, err := http.Get(testUrl + "/")
	testPageStatus(resp, err, http.StatusOK, t)
}

func newServer(t *testing.T) *httptest.Server {
	initDB(t)
	server := httptest.NewServer(Handler())
	testUrl = server.URL
	return server
}

func testPageStatus(resp *http.Response, err error, expectedCode int, t *testing.T) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error", err)
	}
	if resp.StatusCode != expectedCode {
		t.Fatalf("wanted %v, got %v. Response was %v", expectedCode, resp.StatusCode, resp)
	}
}

func TestPledgeWithoutLogin(t *testing.T) {
	server := newServer(t)
	defer server.Close()
	resp, err := http.Get(testUrl + "/pledge")
	testPageStatus(resp, err, http.StatusForbidden, t)
}

func TestThankYouWithoutLogin(t *testing.T) {
	server := newServer(t)
	defer server.Close()
	resp, err := http.Get(testUrl + "/thanks")
	testPageStatus(resp, err, http.StatusForbidden, t)
}

func TestSignin(t *testing.T) {
	server := newServer(t)
	defer server.Close()
	loginUser("test1", t)
}

func loginUser(user string, t *testing.T) http.CookieJar {
	t.Helper()
	resp, err := http.Get(testUrl + "/signin?u=" + user)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("unable to signin, error was %v, response %v", err, resp)
	}
	jar, _ := cookiejar.New(nil)
	url, _ := url.Parse(testUrl)
	jar.SetCookies(url, resp.Cookies())
	return jar
}

func TestPledgeWithLogin(t *testing.T) {
	server := newServer(t)
	defer server.Close()

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

func readResponseBody(resp *http.Response, t *testing.T) string {
	t.Helper()
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

func testStrings(body string, expectedStrings []string, t *testing.T) {
	t.Helper()
	for _, s := range expectedStrings {
		if !strings.Contains(body, s) {
			t.Errorf("body did not contain '%s'", s)
		}
	}
	if t.Failed() {
		t.Errorf("body was %s", body)
	}
}

func TestHappyPathPledge(t *testing.T) {
	server := newServer(t)
	defer server.Close()
	cookie := loginUser("test1", t)
	client := http.Client{Jar: cookie}
	data := url.Values{"item": {"1"}}
	resp, err := client.PostForm(testUrl+"/pledge", data)
	testPageStatus(resp, err, http.StatusOK, t)

	expected := []string{
		"boycott of phone xyz by testco1",
		"Signed in as user1 full",
	}
	body := readResponseBody(resp, t)
	testStrings(body, expected, t)
}

func TestBadNewItems(t *testing.T) {
	server := newServer(t)
	defer server.Close()
	cookie := loginUser("test1", t)
	client := http.Client{Jar: cookie}

	values := []url.Values{
		url.Values{"company": {""}, "make": {"newmake"}, "model": {"newmodel"}},
		url.Values{"company": {"newco"}, "make": {""}, "model": {"newmodel"}},
		url.Values{"company": {""}, "make": {"newmake"}, "model": {""}},
		url.Values{"make": {"newmake"}, "model": {"bla"}},
	}

	for _, d := range values {
		resp, err := client.PostForm(testUrl+"/newitem", d)
		testPageStatus(resp, err, http.StatusBadRequest, t)
	}
	items, err := apis.NewItem.ListNewItems()
	if err != nil {
		t.Fatal(err)
	}
	if 0 != len(items) {
		t.Fatalf("expected 0 new items, got %d", len(items))
	}
}

func TestBlockAccessToItemListing(t *testing.T) {
	server := newServer(t)
	defer server.Close()
	cookie := loginUser("test1", t)
	client := http.Client{Jar: cookie}

	resp, err := client.Get(testUrl + "/admin/newitems")
	testPageStatus(resp, err, http.StatusForbidden, t)
}

func TestListNewItems(t *testing.T) {
	server := newServer(t)
	defer server.Close()
	cookie := loginUser("test1", t)
	client := http.Client{Jar: cookie}
	loadNewItems(client, t)

	testNewItemsPage(newItemPostData, t)
}

var newItemPostData = []url.Values{
	url.Values{"company": {"newco1"}, "make": {"newmake1"}, "model": {"newmodel1"}},
	url.Values{"company": {"newco2"}, "make": {"newmake2"}, "model": {"newmodel2"}},
	url.Values{"company": {"newco3"}, "make": {"newmake3"}, "model": {"newmodel3"}},
	url.Values{"company": {"newco4"}, "make": {"newmake4"}, "model": {"newmodel4"}},
}

func testNewItemsPage(newItems []url.Values, t *testing.T) {

	type newItemHtml struct {
		Id          string
		IsPledge    string
		UserID      string
		UserCompany string
		UserMake    string
		UserModel   string
		UserValue   string
	}

	t.Helper()
	cookie := loginUser("admin", t)
	client := http.Client{Jar: cookie}
	resp, err := client.Get(testUrl + "/admin/newitems")
	testPageStatus(resp, err, http.StatusOK, t)

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		t.Fatal(err)
	}

	rows := doc.Find(".items-table .item-entry")
	if rows.Size() != len(newItems) {
		t.Fatalf("expected %d rows in item listing, got %d", len(newItems), rows.Size())
	}
	rows.Each(func(i int, s *goquery.Selection) {
		h := newItemHtml{}

		fmap := map[string]func(string){
			"id[]":          func(w string) { h.Id = w },
			"isPledge[]":    func(w string) { h.IsPledge = w },
			"userID[]":      func(w string) { h.UserID = w },
			"usercompany[]": func(w string) { h.UserCompany = w },
			"usermake[]":    func(w string) { h.UserMake = w },
			"usermodel[]":   func(w string) { h.UserModel = w },
			"uservalue[]":   func(w string) { h.UserValue = w },
			"add[]":         func(string) {},
			"delete[]":      func(string) {},
		}

		s.Find("input").Each(func(_ int, s *goquery.Selection) {
			n := s.AttrOr("name", "")
			v := s.AttrOr("value", "")
			f, ok := fmap[n]
			if !ok {
				t.Fatalf("unknown input: %s", n)
			}
			f(v)
		})

		v := newItems[i]
		expectedNewItem := newItemHtml{
			Id:          strconv.Itoa(i + 1),
			IsPledge:    "true",
			UserID:      "test1",
			UserCompany: v["company"][0],
			UserMake:    v["make"][0],
			UserModel:   v["model"][0],
			UserValue:   "0",
		}
		if !reflect.DeepEqual(expectedNewItem, h) {
			t.Errorf("expected html form to have %v, got %v", expectedNewItem, h)
		}
	})

}

func loadNewItems(client http.Client, t *testing.T) {
	t.Helper()
	for i, d := range newItemPostData {
		resp, err := client.PostForm(testUrl+"/newitem", d)
		testPageStatus(resp, err, http.StatusOK, t)
		expected := []string{
			fmt.Sprintf("boycott of %s %s by %s", newItemPostData[i]["make"][0], newItemPostData[i]["model"][0], newItemPostData[i]["company"][0]),
			"Signed in as user1 full",
			"new item is under review",
		}
		body := readResponseBody(resp, t)
		testStrings(body, expected, t)
		// Round-tripping of database items is tested in sqldb package, no need
		// to replicate the work here
	}
}

func TestBadNewItemPost(t *testing.T) {
	server := newServer(t)
	defer server.Close()
	cookie := loginUser("admin", t)
	client := http.Client{Jar: cookie}

	post := url.Values{
		"add[]":      []string{"0"}, // row selected
		"isPledge[]": []string{"0"},
		"userID[]":   []string{"test1"}, // user test1
	}
	resp, err := client.PostForm(testUrl+"/admin/newitems", post)
	testPageStatus(resp, err, http.StatusBadRequest, t)

}

func TestNewItemAdminPost(t *testing.T) {
	server := newServer(t)
	defer server.Close()
	cookie := loginUser("test1", t)
	client := http.Client{Jar: cookie}
	loadNewItems(client, t)

	cookie = loginUser("admin", t)
	client = http.Client{Jar: cookie}

	for itemToUse := 3; itemToUse > 0; itemToUse-- {
		post := url.Values{
			"add[]":         []string{"0"}, // row selected
			"id[]":          []string{fmt.Sprintf("%d", itemToUse+1)},
			"isPledge[]":    []string{"0"},
			"userID[]":      []string{"test1"}, // user test1
			"item[]":        []string{"0"},     // add new item
			"company[]":     []string{"0"},     //add new company
			"usercompany[]": newItemPostData[itemToUse]["company"],
			"usermake[]":    newItemPostData[itemToUse]["make"],
			"usermodel[]":   newItemPostData[itemToUse]["model"],
		}
		resp, err := client.PostForm(testUrl+"/admin/newitems", post)
		testPageStatus(resp, err, http.StatusOK, t)

		testNewItemAdded(newItemPostData[itemToUse], t)
		testNewItemsPage(newItemPostData[:itemToUse], t)
	}
}

func testNewItemAdded(item url.Values, t *testing.T) {
	t.Helper()
	items, err := apis.Item.ListItems()
	if err != nil {
		t.Fatal(err)
	}
	make, model, company := item["make"][0], item["model"][0], item["company"][0]
	for _, i := range items {
		if i.Make == make && i.Model == model && i.Company.Name == company {
			return
		}
	}
	t.Fatalf("did not find expected item %s,%s,%s in items db", make, model, company)
}

func TestLimitUserNewItems(t *testing.T) {
}

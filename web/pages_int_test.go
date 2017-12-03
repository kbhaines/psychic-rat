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

func TestNewItems(t *testing.T) {
	server := newServer(t)
	defer server.Close()
	cookie := loginUser("test1", t)
	client := http.Client{Jar: cookie}

	loadNewItems(client, t)

	items, err := apis.NewItem.ListNewItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(newItemPostData) != len(items) {
		t.Fatalf("expected %d new items, got %d", len(newItemPostData), len(items))
	}
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

var newItemPostData = []url.Values{
	url.Values{"company": {"newco1"}, "make": {"newmake1"}, "model": {"newmodel1"}},
	url.Values{"company": {"newco2"}, "make": {"newmake2"}, "model": {"newmodel2"}},
	url.Values{"company": {"newco3"}, "make": {"newmake3"}, "model": {"newmodel3"}},
	url.Values{"company": {"newco4"}, "make": {"newmake4"}, "model": {"newmodel4"}},
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

	cookie = loginUser("admin", t)
	client = http.Client{Jar: cookie}
	resp, err := client.Get(testUrl + "/admin/newitems")
	testPageStatus(resp, err, http.StatusOK, t)

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		t.Fatal(err)
	}

	rows := doc.Find(".items-table .item-entry")
	if rows.Size() != len(newItemPostData) {
		t.Errorf("expected %d rows in results, got %d", len(newItemPostData), rows.Size())
	}
	rows.Each(func(i int, s *goquery.Selection) {
		h := newItemHtml{}
		s.Find("input").Each(func(_ int, s *goquery.Selection) {
			n := s.AttrOr("name", "")
			v := s.AttrOr("value", "")
			switch n {
			case "id[]":
				h.Id = v
			case "isPledge[]":
				h.IsPledge = v
			case "userID[]":
				h.UserID = v
			case "usercompany[]":
				h.UserCompany = v
			case "usermake[]":
				h.UserMake = v
			case "usermodel[]":
				h.UserModel = v
			case "uservalue[]":
				h.UserValue = v
			case "add[]":
			case "delete[]":
				break
			default:
				t.Fatalf("unknown input: %s", n)
			}
		})

		v := newItemPostData[i]
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

type newItemHtml struct {
	Id          string
	IsPledge    string
	UserID      string
	UserCompany string
	UserMake    string
	UserModel   string
	UserValue   string
}

func TestItemAdminSubmission(t *testing.T) {
	server := newServer(t)
	defer server.Close()
	cookie := loginUser("test1", t)
	client := http.Client{Jar: cookie}
	loadNewItems(client, t)

	cookie = loginUser("admin", t)
	client = http.Client{Jar: cookie}

	post := url.Values{
		"add[]":         []string{"0"},
		"id[]":          []string{"1"},
		"isPledge[]":    []string{"0"},
		"userID[]":      []string{"test1"},
		"item[]":        []string{"0"},
		"company[]":     []string{"0"},
		"usercompany[]": newItemPostData[0]["company"],
		"usermake[]":    newItemPostData[0]["make"],
		"usermodel[]":   newItemPostData[0]["model"],
	}
	resp, err := client.PostForm(testUrl+"/admin/newitems", post)
	testPageStatus(resp, err, http.StatusOK, t)
}

func TestLimitUserNewItems(t *testing.T) {
}

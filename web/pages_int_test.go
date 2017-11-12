package web

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"psychic-rat/mdl"
	"strings"
	"testing"
)

// Todo : move mock data into tests

var testUrl = "http://localhost:8080"

var testUsers = []mdl.User{
	mdl.User{Id: "test1", Fullname: "user1 full", FirstName: "user1", Email: "user1@user.com"},
	mdl.User{Id: "test2", Fullname: "user2 full", FirstName: "user2", Email: "user2@user.com"},
	mdl.User{Id: "test3", Fullname: "user3 full", FirstName: "user3", Email: "user3@user.com"},
	mdl.User{Id: "test4", Fullname: "user4 full", FirstName: "user4", Email: "user4@user.com"},
	mdl.User{Id: "test5", Fullname: "user5 full", FirstName: "user5", Email: "user5@user.com"},
	mdl.User{Id: "test6", Fullname: "user6 full", FirstName: "user6", Email: "user6@user.com"},
	mdl.User{Id: "test7", Fullname: "user7 full", FirstName: "user7", Email: "user7@user.com"},
	mdl.User{Id: "test8", Fullname: "user8 full", FirstName: "user8", Email: "user8@user.com"},
	mdl.User{Id: "test9", Fullname: "user9 full", FirstName: "user9", Email: "user9@user.com"},
	mdl.User{Id: "test10", Fullname: "user10 full", FirstName: "user10", Email: "user10@user.com"},
}

func initUsers(t *testing.T) {
	t.Helper()
	for _, u := range testUsers {
		err := apis.User.CreateUser(u)
		if err != nil {
			t.Fatal(err)
		}
	}
}
func TestHomePage(t *testing.T) {
	server := newServer()
	defer server.Close()
	resp, err := http.Get(testUrl + "/")
	testPageStatus(resp, err, http.StatusOK, t)
}

func newServer() *httptest.Server {
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
	server := newServer()
	defer server.Close()
	resp, err := http.Get(testUrl + "/pledge")
	testPageStatus(resp, err, http.StatusForbidden, t)
}

func TestThankYouWithoutLogin(t *testing.T) {
	server := newServer()
	defer server.Close()
	resp, err := http.Get(testUrl + "/thanks")
	testPageStatus(resp, err, http.StatusForbidden, t)
}

func TestSignin(t *testing.T) {
	server := newServer()
	defer server.Close()
	initUsers(t)
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
	server := newServer()
	defer server.Close()
	initUsers(t)

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
	server := newServer()
	defer server.Close()
	cookie := loginUser("test1", t)
	client := http.Client{Jar: cookie}
	data := url.Values{"item": {"1"}}
	resp, err := client.PostForm(testUrl+"/pledge", data)
	testPageStatus(resp, err, http.StatusOK, t)

	expected := []string{
		"boycott of phone abc by bigco1",
		"Signed in as Kevin",
	}
	body := readResponseBody(resp, t)
	testStrings(body, expected, t)
}

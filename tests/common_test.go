package tests

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"psychic-rat/auth/basic"
	"psychic-rat/sqldb"
	"psychic-rat/web"
	"psychic-rat/web/admin"
	"psychic-rat/web/pub"
	"psychic-rat/web/tmpl"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var testUrl string

func newServer(t *testing.T) (*httptest.Server, *sqldb.DB) {
	server := httptest.NewServer(web.Handler())
	testUrl = server.URL
	db := initDB(t)
	renderer := tmpl.NewRenderer("../res/tmpl", false)
	authHandler := basic.NewUserHandler()
	web.Init(authHandler)
	pub.Init(db, db, db, authHandler, renderer)
	admin.Init(db, db, db, db, authHandler, renderer)
	return server, db
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

func cleanUp(server *httptest.Server, db *sqldb.DB) {
	server.Close()
	db.Close()
}

func loginUser(user string, t *testing.T) http.CookieJar {
	t.Helper()
	req, err := http.NewRequest("GET", testUrl+"/callback?u="+user, nil)
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil || (resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusSeeOther) {
		t.Fatalf("unable to signin, error was %v, response %v", err, resp)
	}
	jar, _ := cookiejar.New(nil)
	url, _ := url.Parse(testUrl)
	jar.SetCookies(url, resp.Cookies())
	return jar
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

func execAuthdRequest(username, method, url string, postValues url.Values) (*http.Response, error) {
	var req *http.Request
	var err error
	if postValues == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, strings.NewReader(postValues.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, "")
	client := http.Client{}
	return client.Do(req)
}

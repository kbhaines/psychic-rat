package tests

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"psychic-rat/auth"
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

const callback = "/callback?p=basic"

var testURL string

type mockRateLimit struct{}

func (mr *mockRateLimit) CheckLimit(*http.Request) error { return nil }

type mockHuman struct{}

func (m *mockHuman) IsHuman(f url.Values) error { return nil }

func newServer(t *testing.T) (*httptest.Server, *sqldb.DB) {
	server := httptest.NewServer(web.Handler())
	testURL = server.URL
	db := initDB(t)
	renderer := tmpl.NewRenderer("../res/tmpl", false)
	authHandler := auth.NewUserHandler()
	authProviders := map[string]auth.AuthHandler{
		"basic": basic.New(testURL + callback),
	}
	auth.Init(db, authProviders)
	web.Init(authHandler, &mockRateLimit{})
	pub.Init(db, db, db, authHandler, renderer, &mockHuman{})
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

func getAuthdClient(user string, t *testing.T) *http.Client {
	t.Helper()
	req, err := http.NewRequest("GET", testURL+callback, nil)
	req.SetBasicAuth(user, "")
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil || (resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusSeeOther) {
		t.Fatalf("unable to signin, error was %v, response %v", err, resp)
	}
	jar, _ := cookiejar.New(nil)
	url, _ := url.Parse(testURL)
	jar.SetCookies(url, resp.Cookies())
	return &http.Client{Jar: jar}
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

func getCSRFToken(client *http.Client, srcURL string, t *testing.T) string {
	t.Helper()
	resp, err := client.Get(srcURL)
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

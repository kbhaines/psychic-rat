// +build integration

package main

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"testing"
)

var testUrl = "http://localhost:8080"

func TestHomePage(t *testing.T) {
	resp, err := http.Get(testUrl + "/")
	testPageStatus(resp, err, http.StatusOK, t)
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
	resp, err := http.Get(testUrl + "/pledge")
	testPageStatus(resp, err, http.StatusForbidden, t)
}

func TestThankYouWithoutLogin(t *testing.T) {
	resp, err := http.Get(testUrl + "/thanks")
	testPageStatus(resp, err, http.StatusForbidden, t)
}

func TestSignin(t *testing.T) {
	loginUser("testuser1", t)
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
	cookie := loginUser("testuser1", t)
	client := http.Client{Jar: cookie}
	resp, err := client.Get(testUrl + "/pledge")
	testPageStatus(resp, err, http.StatusOK, t)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	strBody := string(body)
	if !strings.Contains(strBody, "Kevin") {
		t.Error("body does not contain greeting for user")
	}
	if !strings.Contains(strBody, "<select ") {
		t.Error("body does not contain <select ")
	}
	if !strings.Contains(strBody, "<input type=\"submit\"") {
		t.Error("body does not contain <input> for submit button")
	}
	// Todo : move mock data
	// Todo: move to using handler-based testing
}

func TestHappyPathPledge(t *testing.T) {
	cookie := loginUser("testuser1", t)
	client := http.Client{Jar: cookie}
	data := url.Values{"item": {"2"}}
	resp, err := client.PostForm(testUrl+"/pledge", data)
	testPageStatus(resp, err, http.StatusOK, t)
}

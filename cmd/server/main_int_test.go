// +build integration

package main

import (
	"net/http"
	"testing"
)

var testUrl = "http://localhost:8080"

func TestHomePage(t *testing.T) {
	testPageStatus("/", http.StatusOK, t)
}

func TestPledgeWithoutLogin(t *testing.T) {
	testPageStatus("/pledge", http.StatusForbidden, t)
}

func TestThankYouWithoutLogin(t *testing.T) {
	testPageStatus("/thankyou", http.StatusForbidden, t)
}

func testPageStatus(page string, expectedCode int, t *testing.T) {
	t.Helper()
	resp, err := http.Get(testUrl + page)
	if err != nil {
		t.Fatalf("unexpected error accessing %v", err, page)
	}
	if resp.StatusCode != expectedCode {
		t.Fatalf("wanted %v, got %v : was able to access %v page without login", expectedCode, resp.StatusCode, page)
	}
}

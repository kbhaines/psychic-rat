package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"psychic-rat/sess"
	"testing"
)

func TestCSRFBlock(t *testing.T) {
	values := url.Values{"item": {"1"}}
	req := &http.Request{Method: "POST", PostForm: values}
	writer := httptest.NewRecorder()
	called := false
	csrfFunc := csrfProtect(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	csrfFunc(writer, req)
	if called {
		t.Fatalf("protection failed to stop handler call")
	}
	if writer.Result().StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403, got %s", writer.Result().Status)
	}
}

func TestCSRFPass(t *testing.T) {
	writer := httptest.NewRecorder()
	req := &http.Request{Method: "POST"}
	token, err := sess.NewSessionStore(req).SetCSRF(writer)
	if err != nil {
		t.Fatal(err)
	}
	req.PostForm = url.Values{"item": {"1"}, "csrf": {token}}
	called := false
	csrfFunc := csrfProtect(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	csrfFunc(writer, req)
	if !called {
		t.Fatalf("handler not called")
	}
}

func TestCSRFIgnore(t *testing.T) {
	writer := httptest.NewRecorder()
	req := &http.Request{Method: "GET"}
	called := false
	csrfFunc := csrfProtect(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	csrfFunc(writer, req)
	if !called {
		t.Fatalf("handler not called")
	}
}

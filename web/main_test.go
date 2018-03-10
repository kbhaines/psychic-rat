package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"psychic-rat/types"
	"testing"
)

type mockUserHandler struct {
	csrfResult   bool
	tokenChecked string
}

func (m *mockUserHandler) GetLoggedInUser(*http.Request) (*types.User, error) {
	panic("not implemented")
}

func (m *mockUserHandler) VerifyUserCSRF(r *http.Request, token string) error {
	m.tokenChecked = token
	if m.csrfResult {
		return nil
	}
	return fmt.Errorf("CSRF failed")
}

type mockRateLimit struct{}

func (mr *mockRateLimit) CheckLimit(*http.Request) error { return nil }

func TestCSRFBlock(t *testing.T) {
	values := url.Values{"item": {"1"}}
	req := &http.Request{Method: "POST", PostForm: values}
	writer := httptest.NewRecorder()
	called := false
	Init(&mockUserHandler{csrfResult: false}, &mockRateLimit{})
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
	token := "1234"
	req.PostForm = url.Values{"item": {"1"}, "csrf": {token}}
	mock := &mockUserHandler{csrfResult: true}
	Init(mock, &mockRateLimit{})
	called := false
	csrfFunc := csrfProtect(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	csrfFunc(writer, req)
	if !called {
		t.Fatalf("handler not called")
	}
	if mock.tokenChecked != token {
		t.Fatalf("expected token = %s in call to VerifyUserCSRF, got %s", token, mock.tokenChecked)
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

package auth0

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func handler() http.Handler {
	hmux := http.NewServeMux()
	hmux.HandleFunc("/callback", CallbackHandler)
	return hmux
}

func mockAuth0() http.Handler {
	hmux := http.NewServeMux()
	hmux.HandleFunc("/oauth/token", tokenHandler)
	hmux.HandleFunc("/userinfo", loggingHandler)
	hmux.HandleFunc("/authorize", loggingHandler)
	hmux.HandleFunc("/*", loggingHandler)
	return hmux
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("*r = %+v\n", *r)
}

func loggingHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("*r = %+v\n", *r)
}

func TestCallbackHandler(t *testing.T) {
	server := httptest.NewServer(handler())
	defer server.Close()

	auth0Server := httptest.NewServer(mockAuth0())
	defer auth0Server.Close()

	Init(nil)
	os.Setenv("AUTH0_CLIENT_ID", "client_id")
	os.Setenv("AUTH0_CLIENT_SECRET", "secret")
	os.Setenv("AUTH0_CALLBACK_URL", server.URL+"/callback")
	os.Setenv("AUTH0_DOMAIN", "localhost")
	resp, err := http.Get(server.URL + "/callback?code=12345")
	if err != nil {
		t.Fatalf("error in getting callback: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusSeeOther {
		//body, _ := ioutil.ReadAll(resp.Body)
		//t.Fatalf("did not get redirected; status was %v response was %v, body was %v", resp.Status, resp, string(body))
	}
	expectLocation := "some"
	actualLocation := resp.Header.Get("location")
	if actualLocation != expectLocation {
		//t.Fatalf("expected redirect to %v, got %v", expectLocation, actualLocation)
	}
}

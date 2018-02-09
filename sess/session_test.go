package sess

import (
	"net/http"
	"net/http/httptest"
	"psychic-rat/types"
	"reflect"
	"testing"
)

func TestSaveAndGet(t *testing.T) {
	r := http.Request{}
	w := httptest.NewRecorder()
	sess := NewSessionStore(&r)
	user := &types.User{ID: "test1", Email: "me@test1.com"}

	saveUser(user, sess, w, t)
	gotUser := getUser(sess, t)

	if !reflect.DeepEqual(*user, *gotUser) {
		t.Fatalf("expected user %v, got %v", *user, *gotUser)
	}

}

func saveUser(user *types.User, sess *SessionStore, w http.ResponseWriter, t *testing.T) {
	t.Helper()
	if err := sess.Save(user, w); err != nil {
		t.Fatal(err)
	}
}

func getUser(sess *SessionStore, t *testing.T) *types.User {
	t.Helper()
	user, err := sess.Get()
	if err != nil {
		t.Fatal(err)
	}
	return user
}

func TestSaveAndClear(t *testing.T) {
	r := http.Request{}
	w := httptest.NewRecorder()
	sess := NewSessionStore(&r)
	user := &types.User{ID: "test1", Email: "me@test1.com"}

	if err := sess.Save(user, w); err != nil {
		t.Fatal(err)
	}

	sess.Save(nil, w)
	gotUser := getUser(sess, t)
	if gotUser != nil {
		t.Fatalf("expected nil user, got %v", *gotUser)
	}
}

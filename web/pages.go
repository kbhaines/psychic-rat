package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"psychic-rat/sess"
	"psychic-rat/types"
	"psychic-rat/web/dispatch"
	"psychic-rat/web/tmpl"
	"strconv"

	"github.com/gorilla/sessions"
)

// TODO: split into subpackages; separate user & admin stuff for starters.

type (
	authInfo struct {
		Auth0ClientId    string
		Auth0CallbackURL string
		Auth0Domain      string
	}

	// pageVariables holds data for the templates. Stuffed into one struct for now.
	pageVariables struct {
		authInfo
		Items     []types.Item
		User      types.User
		NewItems  []types.NewItem
		Companies []types.Company
	}
)

var (
	// function variables, allows us to swap out for mocks for easier testing
	renderPage     = tmpl.RenderTemplate
	isUserLoggedIn = isUserLoggedInSession

	// TODO: Env var
	auth0Store = sessions.NewCookieStore([]byte("something-very-secret"))
)

func HomePageHandler(writer http.ResponseWriter, request *http.Request) {
	selector := dispatch.MethodSelector{
		"GET": func(writer http.ResponseWriter, request *http.Request) {
			vars := (&pageVariables{}).withSessionVars(request)
			renderPage(writer, "home.html.tmpl", vars)
		},
	}
	dispatch.ExecHandlerForMethod(selector, writer, request)
}

func SignInPageHandler(writer http.ResponseWriter, request *http.Request) {
	method := signInSimple
	if flags.enableAuth0 {
		method = signInAuth0
	}
	dispatch.ExecHandlerForMethod(dispatch.MethodSelector{"GET": method}, writer, request)
}

func signInSimple(writer http.ResponseWriter, request *http.Request) {
	s := sess.NewSessionStore(request, writer)
	user, err := s.Get()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		log.Print("no user, attempting auth")
		if err := authUser(request, s); err != nil {
			log.Print(err)
			http.Error(writer, "authentication failed", http.StatusForbidden)
			return
		}
	}

	vars := (&pageVariables{}).withSessionVars(request)
	renderPage(writer, "signin.html.tmpl", vars)
}

func authUser(request *http.Request, session *sess.SessionStore) error {
	if err := request.ParseForm(); err != nil {
		return err
	}

	userId := request.FormValue("u")
	if userId == "" {
		return fmt.Errorf("userId not specified")
	}

	user, err := apis.User.GetUser(userId)
	if err != nil {
		return fmt.Errorf("can't get user by id %v : %v", userId, err)
	}
	return session.Save(*user)
}

func signInAuth0(writer http.ResponseWriter, request *http.Request) {
	vars := (&pageVariables{}).withAuth0Vars()
	log.Printf("vars = %+v\n", vars)
	renderPage(writer, "signin-auth0.html.tmpl", vars)
}

func userLoginRequired(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isUserLoggedIn(r) {
			log.Print("user not logged in")
			http.Error(w, "", http.StatusForbidden)
			return
		}
		h(w, r)
	}
}

func isUserLoggedInSession(request *http.Request) bool {
	// TODO: nil is a smell. StoreReader/Writer interfaces.
	s := sess.NewSessionStore(request, nil)
	user, err := s.Get()
	if err != nil {
		log.Print(err)
		return false
	}
	return user != nil
}

func PledgePageHandler(writer http.ResponseWriter, request *http.Request) {
	selector := dispatch.MethodSelector{
		"GET":  userLoginRequired(pledgeGetHandler),
		"POST": userLoginRequired(pledgePostHandler),
	}
	dispatch.ExecHandlerForMethod(selector, writer, request)
}

func pledgeGetHandler(writer http.ResponseWriter, request *http.Request) {
	report, err := apis.Item.ListItems()
	if err != nil {
		log.Print(err)
		http.Error(writer, "", http.StatusInternalServerError)
		return
	}
	vars := &pageVariables{Items: report}
	vars = vars.withSessionVars(request)
	renderPage(writer, "pledge.html.tmpl", vars)
}

func pledgePostHandler(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		log.Print(err)
		http.Error(writer, "", http.StatusInternalServerError)
		return
	}

	itemId64, err := strconv.ParseInt(request.FormValue("item"), 10, 32)
	if err != nil {
		http.Error(writer, "", http.StatusBadRequest)
		return
	}
	itemId := int(itemId64)

	item, err := apis.Item.GetItem(itemId)
	if err != nil {
		log.Printf("error looking up item %v : %v", itemId, err)
		http.Error(writer, "", http.StatusBadRequest)
		return
	}

	// TODO ignoring a couple of errors
	s := sess.NewSessionStore(request, writer)
	user, _ := s.Get()
	userId := user.ID

	log.Printf("pledge item %v from user %v", itemId, userId)
	plId, err := apis.Pledge.AddPledge(itemId, userId)
	if err != nil {
		log.Print("unable to pledge : ", err)
		return
	}
	log.Printf("pledge %v created", plId)

	vars := &pageVariables{
		User:  *user,
		Items: []types.Item{item},
	}
	renderPage(writer, "pledge-post.html.tmpl", vars)
}

func ThanksPageHandler(writer http.ResponseWriter, request *http.Request) {
	selector := dispatch.MethodSelector{
		"GET": userLoginRequired(func(writer http.ResponseWriter, request *http.Request) {
			vars := (&pageVariables{}).withSessionVars(request)
			renderPage(writer, "thanks.html.tmpl", vars)
		}),
	}
	dispatch.ExecHandlerForMethod(selector, writer, request)
}

func NewItemHandler(w http.ResponseWriter, r *http.Request) {
	selector := dispatch.MethodSelector{
		"POST": userLoginRequired(newItemPostHandler),
	}
	dispatch.ExecHandlerForMethod(selector, w, r)
}

func newItemPostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Print(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	// TODO ignoring a couple of errors
	s := sess.NewSessionStore(r, w)
	user, _ := s.Get()
	userId := user.ID

	company := r.FormValue("company")
	make := r.FormValue("make")
	model := r.FormValue("model")
	if company == "" || model == "" || make == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	//value, err := strconv.ParseFloat(r.FormValue("value"), 10)
	//if err != nil {
	//	log.Print(err)
	//	http.Error(w, "", http.StatusBadRequest)
	//	return
	//}

	newItem := types.NewItem{UserID: userId, IsPledge: true, Make: make, Model: model, Company: company}
	_, err := apis.NewItem.AddNewItem(newItem)
	if err != nil {
		log.Printf("unable to add new item %v:", newItem, err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	item := types.Item{Make: newItem.Make, Model: newItem.Model, Company: types.Company{Name: company}}
	vars := &pageVariables{
		User:  *user,
		Items: []types.Item{item},
	}
	renderPage(w, "pledge-post-new-item.html.tmpl", vars)
}

func (pv *pageVariables) withSessionVars(r *http.Request) *pageVariables {
	// TODO: nil is a smell. StoreReader/Writer interfaces.
	s := sess.NewSessionStore(r, nil)
	user, err := s.Get()
	if err != nil {
		// todo - return error?
		log.Fatal(err)
	}
	if user != nil {
		pv.User = *user
	}
	return pv
}

func (pv *pageVariables) withAuth0Vars() *pageVariables {
	pv.Auth0Domain = os.Getenv("AUTH0_DOMAIN")
	pv.Auth0CallbackURL = os.Getenv("AUTH0_CALLBACK_URL")
	pv.Auth0ClientId = os.Getenv("AUTH0_CLIENT_ID")
	return pv
}

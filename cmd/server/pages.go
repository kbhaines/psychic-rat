package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"psychic-rat/mdl"
	"psychic-rat/sess"
	"psychic-rat/types"

	"github.com/gorilla/sessions"
)

type (
	Auth0 struct {
		Auth0ClientId    string
		Auth0CallbackURL string
		Auth0Domain      string
	}

	pageVariables struct {
		Auth0
		Items []types.ItemElement
		User  mdl.UserRecord
	}

	renderFunc     func(writer http.ResponseWriter, templateName string, vars *pageVariables)
	handlerFunc    func(http.ResponseWriter, *http.Request)
	methodSelector map[string]handlerFunc
)

var (
	renderPage     = renderPageUsingTemplate
	isUserLoggedIn = isUserLoggedInSession
	auth0Store     = sessions.NewCookieStore([]byte("something-very-secret"))
	logDbg         = log.New(os.Stderr, "DBG:", 0).Print
	logDbgf        = log.New(os.Stderr, "DBG:", 0).Printf
)

func renderPageUsingTemplate(writer http.ResponseWriter, templateName string, variables *pageVariables) {
	tpt := template.Must(template.New(templateName).ParseFiles(templateName, "header.html.tmpl", "footer.html.tmpl", "navi.html.tmpl"))
	tpt.Execute(writer, variables)
}

func HomePageHandler(writer http.ResponseWriter, request *http.Request) {
	selector := methodSelector{
		"GET": func(writer http.ResponseWriter, request *http.Request) {
			vars := (&pageVariables{}).withSessionVars(request)
			renderPage(writer, "home.html.tmpl", vars)
		},
	}
	execHandlerForMethod(selector, writer, request)
}

func (pv *pageVariables) withSessionVars(r *http.Request) *pageVariables {
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

func execHandlerForMethod(selector methodSelector, writer http.ResponseWriter, request *http.Request) {
	f, exists := selector[request.Method]
	if !exists {
		http.Error(writer, "", http.StatusMethodNotAllowed)
		return
	}
	f(writer, request)
}

func SignInPageHandler(writer http.ResponseWriter, request *http.Request) {
	selector := methodSelector{
		"GET": signInSimple,
	}
	execHandlerForMethod(selector, writer, request)
}

func signInSimple(writer http.ResponseWriter, request *http.Request) {
	s := sess.NewSessionStore(request, writer)
	user, err := s.Get()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if user != nil {
		logDbg("no user, attempting auth")
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

	user, err := apis.User.GetById(mdl.Id(userId))
	if err != nil {
		return fmt.Errorf("can't get user by id %v : %v", userId, err)
	}
	return session.Save(*user)

}
func signInAuth0(writer http.ResponseWriter, request *http.Request) {
	vars := (&pageVariables{}).withAuth0Vars()
	log.Printf("vars = %+v\n", vars)
	renderPage(writer, "signin.html.tmpl", vars)
}

func userLoginRequired(h handlerFunc) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isUserLoggedIn(r) {
			log.Print("user not logged in")
			http.Error(w, "", http.StatusForbidden)
			return
		}
		h(w, r)
	}
}

func PledgePageHandler(writer http.ResponseWriter, request *http.Request) {
	selector := methodSelector{
		"GET":  userLoginRequired(pledgeGetHandler),
		"POST": userLoginRequired(pledgePostHandler),
	}
	execHandlerForMethod(selector, writer, request)
}

func isUserLoggedInSession(request *http.Request) bool {
	s := sess.NewSessionStore(request, nil)
	user, err := s.Get()
	if err != nil {
		log.Print(err)
		return false
	}
	return user != nil
}

func pledgeGetHandler(writer http.ResponseWriter, request *http.Request) {
	report, err := apis.Item.ListItems()
	if err != nil {
		log.Print(err)
		http.Error(writer, "", http.StatusInternalServerError)
		return
	}
	vars := &pageVariables{Items: report.Items}
	vars = vars.withSessionVars(request)
	renderPage(writer, "pledge.html.tmpl", vars)
}

func pledgePostHandler(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		log.Print(err)
		http.Error(writer, "", http.StatusInternalServerError)
		return
	}
	logDbgf("request = %+v\n", request)
	logDbgf("request.Form= %+v\n", request.Form)

	itemId := mdl.Id(request.FormValue("item"))
	if itemId == "" {
		http.Error(writer, "", http.StatusBadRequest)
		return
	}

	if _, err := apis.Item.GetById(itemId); err != nil {
		log.Printf("error looking up item %v : %v", itemId, err)
		http.Error(writer, "", http.StatusBadRequest)
		return
	}

	// TODO ignoring a couple of errors
	s := sess.NewSessionStore(request, writer)
	user, _ := s.Get()
	userId := user.Id

	log.Printf("pledge item %v from user %v", itemId, userId)
	plId, err := apis.Pledge.NewPledge(itemId, userId)
	if err != nil {
		log.Print("unable to pledge : ", err)
		return
	}
	log.Printf("pledge %v created", plId)
	vars := struct {
		Username string
		PledgeId string
	}{string(userId), string(plId)}
	_ = vars
	//renderPage(writer, "thanks.html.tmpl", vars)
	return
}

func ThanksPageHandler(writer http.ResponseWriter, request *http.Request) {
	selector := methodSelector{
		"GET": userLoginRequired(func(writer http.ResponseWriter, request *http.Request) {
			vars := (&pageVariables{}).withSessionVars(request)
			renderPage(writer, "thanks.html.tmpl", vars)
		}),
	}
	execHandlerForMethod(selector, writer, request)
}

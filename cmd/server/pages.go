package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"psychic-rat/mdl"
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
	store          = sessions.NewCookieStore([]byte("something-very-secret"))
	auth0Store     = sessions.NewCookieStore([]byte("something-very-secret"))
	logDbg         = log.New(os.Stderr, "DBG:", 0).Print
	logDbgf        = log.New(os.Stderr, "DBG:", 0).Printf
)

func renderPageUsingTemplate(writer http.ResponseWriter, templateName string, variables *pageVariables) {
	tpt := template.Must(template.New(templateName).ParseFiles(templateName, "header.html.tmpl", "footer.html.tmpl"))
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
	session, err := auth0Store.Get(r, "auth-session")
	if err != nil {
		return pv
	}
	user, ok := session.Values["userRecord"]
	if !ok {
		return pv
	}
	pv.User = user.(mdl.UserRecord)
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
		"GET": signInAuth0,
	}
	execHandlerForMethod(selector, writer, request)
}

func signInSimple(writer http.ResponseWriter, request *http.Request) {
	session, err := store.Get(request, "session")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	val := session.Values["userId"]
	if _, ok := val.(string); !ok {
		logDbg("no user, attempting auth")
		if err := authUser(session, request); err != nil {
			log.Print(err)
			http.Error(writer, "authentication failed", http.StatusForbidden)
			return
		}
	}
	userId := session.Values["userId"].(string)
	log.Printf("session user is %v", userId)

	if err := session.Save(request, writer); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	vars := (&pageVariables{}).withSessionVars(request)
	renderPage(writer, "signin.html.tmpl", vars)
}

func authUser(session *sessions.Session, request *http.Request) error {
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
	session.Values["userId"] = userId
	session.Values["userEmail"] = user.Email
	return nil

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
	session, err := store.Get(request, "session")
	if err != nil {
		log.Print(err)
		return false
	}
	_, ok := session.Values["userRecord"].(mdl.UserRecord)
	return ok
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

	session, _ := store.Get(request, "session")
	userId := mdl.Id(session.Values["userId"].(string))

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

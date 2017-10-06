package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"psychic-rat/api"
	"psychic-rat/mdl"

	"github.com/gorilla/sessions"
)

type (
	pageVariables struct {
		Username string
		Items    []api.ItemElement
	}
	renderFunc     func(writer http.ResponseWriter, templateName string, variables interface{})
	handlerFunc    func(http.ResponseWriter, *http.Request)
	methodSelector map[string]handlerFunc
)

var (
	renderPage = renderPageUsingTemplate
	store      = sessions.NewCookieStore([]byte("something-very-secret"))
	logDbg     = log.New(os.Stderr, "DBG:", 0).Print
	logDbgf    = log.New(os.Stderr, "DBG:", 0).Printf
)

func renderPageUsingTemplate(writer http.ResponseWriter, templateName string, variables interface{}) {
	tpt := template.Must(template.New(templateName).ParseFiles(templateName, "header.html.tmpl", "footer.html.tmpl"))
	tpt.Execute(writer, variables)
}

func HomePageHandler(writer http.ResponseWriter, request *http.Request) {
	selector := methodSelector{
		"GET": func(writer http.ResponseWriter, request *http.Request) {
			vars := pageVariables{Username: "Kevin"}
			renderPage(writer, "home.html.tmpl", vars)
		},
	}
	execHandlerForMethod(selector, writer, request)
}

func execHandlerForMethod(selector methodSelector, writer http.ResponseWriter, request *http.Request) {
	f, exists := selector[request.Method]
	if !exists {
		log.Print("invalid method in request ", request)
		return
	}
	f(writer, request)
}

func SignInPageHandler(writer http.ResponseWriter, request *http.Request) {
	selector := methodSelector{
		"GET": signIn,
	}
	execHandlerForMethod(selector, writer, request)
}

func signIn(writer http.ResponseWriter, request *http.Request) {
	session, err := store.Get(request, "session")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	val := session.Values["userId"]
	if _, ok := val.(string); !ok {
		logDbg("no user, attempting auth")
		if err := authUser(session, request); err != nil {
			log.Printf("auth failed")
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	userId := session.Values["userId"].(string)
	log.Printf("session user is %v", userId)

	if err := session.Save(request, writer); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	vars := pageVariables{Username: session.Values["userEmail"].(string)}
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
		return fmt.Errorf("can't get user %v : %v", userId, err)
	}
	session.Values["userId"] = userId
	session.Values["userEmail"] = user.Email
	return nil
}

func PledgePageHandler(writer http.ResponseWriter, request *http.Request) {
	if !isUserLoggedIn(request) {
		logDbg("no user logged in")
		return
	}

	selector := methodSelector{
		"GET":  pledgeGetHandler,
		"POST": pledgePostHandler,
	}
	execHandlerForMethod(selector, writer, request)
}

func isUserLoggedIn(request *http.Request) bool {
	session, err := store.Get(request, "session")
	if err != nil {
		log.Print("could not retrieve session or create new")
		return false
	}
	_, ok := session.Values["userId"].(string)
	return ok
}

func pledgeGetHandler(writer http.ResponseWriter, request *http.Request) {
	report, err := apis.Item.ListItems()
	if err != nil {
		log.Fatal(err)
	}
	vars := pageVariables{Username: "Kevin", Items: report.Items}
	renderPage(writer, "pledge.html.tmpl", vars)
}

func pledgePostHandler(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		log.Print(err)
		return
	}
	logDbgf("request = %+v\n", request)
	logDbgf("request.Form= %+v\n", request.Form)
	itemId := mdl.Id(request.FormValue("item"))
	if itemId == "" {
		log.Print("item field not passed")
		return
	}
	if _, err := apis.Item.GetById(itemId); err != nil {
		log.Printf("error looking up item %v : %v", itemId, err)
		return
	}
	log.Printf("pledge item %v from user", itemId)
	userId := getUserIdFromRequest(request)
	plId, err := apis.Pledge.NewPledge(itemId, userId)
	if err != nil {
		log.Print("unable to pledge : ", err)
		return
	}
	log.Printf("pledge %v created", plId)
	return
}

func ThanksPageHandler(writer http.ResponseWriter, request *http.Request) {
	selector := methodSelector{
		"GET": func(writer http.ResponseWriter, request *http.Request) {
			vars := pageVariables{Username: "Kevin"}
			renderPage(writer, "thanks.html.tmpl", vars)
		},
	}
	execHandlerForMethod(selector, writer, request)
}

func getUserIdFromRequest(request *http.Request) mdl.Id {
	return mdl.Id("1234")
}

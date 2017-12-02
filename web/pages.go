package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"psychic-rat/mdl"
	"psychic-rat/sess"
	"psychic-rat/types"
	"strconv"

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
		Items     []types.Item
		User      mdl.User
		NewItems  []types.NewItem
		Companies []types.Company
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

func execHandlerForMethod(selector methodSelector, writer http.ResponseWriter, request *http.Request) {
	f, exists := selector[request.Method]
	if !exists {
		http.Error(writer, "", http.StatusMethodNotAllowed)
		return
	}
	f(writer, request)
}

func SignInPageHandler(writer http.ResponseWriter, request *http.Request) {
	method := signInSimple
	if flags.enableAuth0 {
		method = signInAuth0
	}
	execHandlerForMethod(methodSelector{"GET": method}, writer, request)
}

func signInSimple(writer http.ResponseWriter, request *http.Request) {
	s := sess.NewSessionStore(request, writer)
	user, err := s.Get()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
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

func adminLoginRequired(h handlerFunc) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isUserAdmin(r) {
			log.Print("user not admin")
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
	// TODO: nil is a smell. StoreReader/Writer interfaces.
	s := sess.NewSessionStore(request, nil)
	user, err := s.Get()
	if err != nil {
		log.Print(err)
		return false
	}
	return user != nil
}

func isUserAdmin(request *http.Request) bool {
	s := sess.NewSessionStore(request, nil)
	user, err := s.Get()
	if err != nil {
		log.Print(err)
		return false
	}
	return user != nil && user.IsAdmin
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
	userId := user.Id

	log.Printf("pledge item %v from user %v", itemId, userId)
	plId, err := apis.Pledge.NewPledge(itemId, userId)
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
	selector := methodSelector{
		"GET": userLoginRequired(func(writer http.ResponseWriter, request *http.Request) {
			vars := (&pageVariables{}).withSessionVars(request)
			renderPage(writer, "thanks.html.tmpl", vars)
		}),
	}
	execHandlerForMethod(selector, writer, request)
}

func NewItemHandler(w http.ResponseWriter, r *http.Request) {
	selector := methodSelector{
		"GET":  userLoginRequired(newItemListHandler),
		"POST": userLoginRequired(newItemPostHandler),
	}
	execHandlerForMethod(selector, w, r)
}

func newItemListHandler(w http.ResponseWriter, r *http.Request) {

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
	userId := user.Id

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

func AdminItemHandler(w http.ResponseWriter, r *http.Request) {
	selector := methodSelector{
		"GET":  adminLoginRequired(listNewItems),
		"POST": adminLoginRequired(approveNewItems),
	}
	execHandlerForMethod(selector, w, r)
}

func listNewItems(w http.ResponseWriter, r *http.Request) {
	newItems, err := apis.NewItem.ListNewItems()
	if err != nil {
		log.Printf("unable to retrieve new items: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	items, err := apis.Item.ListItems()
	if err != nil {
		log.Printf("unable to retrieve items: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	companies, err := apis.Company.GetCompanies()
	if err != nil {
		log.Printf("unable to retrieve companies: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	vars := &pageVariables{Items: items, NewItems: newItems, Companies: companies}
	renderPage(w, "admin-new-items.html.tmpl", vars)
}

func approveNewItems(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("unable to parse form: %v", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	addItems, ok := r.Form["add[]"]
	if !ok {
		log.Printf("no add items\n")
		log.Printf("r.Form = %+v\n", r.Form)
		return
	}
	for _, rowString := range addItems {
		row, err := strconv.ParseInt(rowString, 10, 32)
		if err != nil {
			log.Printf("unable to parse row Id: %v", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		log.Printf("add row = %+v\n", row)

	}
}

package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
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
	reader := newFormReader(r.Form)
	for reader.next() {
		nip := reader.getNewItemPost()
		if reader.errors() {
			break
		}

		txn := apiTxn{nil, apis}
		var newItem *types.Item
		if nip.ItemID == 0 {
			var company *types.Company
			if nip.CompanyID == 0 {
				company = txn.addCompany(types.Company{Name: nip.UserCompany})
			} else {
				company = txn.getCompany(nip.CompanyID)
			}
			newItem = txn.addItem(types.Item{Company: *company, Make: nip.UserMake, Model: nip.UserModel})
		} else {
			newItem = txn.getItem(nip.ItemID)
		}

		if nip.Pledge {
			txn.addPledge(newItem.Id, nip.UserID)
		}
		txn.deleteNewItem(nip.Id)
		if txn.err != nil {
			log.Printf("unable to complete transaction for new item %d:  %v", nip.Id, err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

	}
	if reader.errors() {
		log.Printf("errors while parsing form line %d: %v", reader.row, err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}

// formReader parses a submitted newItems form POST request
type formReader struct {
	form url.Values
	row  int
	rows []int
	err  []error
}

func newFormReader(form url.Values) *formReader {
	fr := &formReader{form, -1, []int{}, nil}
	adds, ok := form["add[]"]
	if !ok {
		return fr
	}
	for _, str := range adds {
		rowID, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			log.Printf("unable to parse row Id: %v", err)
			return fr
		}
		fr.rows = append(fr.rows, int(rowID))
	}
	return fr
}

func (f *formReader) errors() bool {
	return len(f.err) > 0
}

func (f *formReader) next() bool {
	if f.errors() {
		return false
	}
	f.row++
	return f.row < len(f.rows)
}

func (f *formReader) getNewItemPost() NewItemAdminPost {
	i := NewItemAdminPost{
		Id:          f.getInt("id[]"),
		UserID:      f.getString("userID[]"),
		ItemID:      f.getInt("item[]"),
		CompanyID:   f.getInt("company[]"),
		UserCompany: f.getString("usercompany[]"),
		UserMake:    f.getString("usermake[]"),
		UserModel:   f.getString("usermodel[]"),
		Pledge:      f.getString("isPledge[]") == "1",
	}
	return i
}

func (f *formReader) getString(field string) string {
	v, ok := f.form[field]
	if !ok || !(f.row < len(v)) {
		f.err = append(f.err, fmt.Errorf("%s not found in form (looking up row %d in %d items)", field, f.row, len(v)))
		return ""
	}
	return v[f.row]
}

func (f *formReader) getInt(field string) int {
	val := f.getString(field)
	i, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		f.err = append(f.err, fmt.Errorf("error parsing field %s = %s into int: %v", field, val, err))
		return 0
	}
	return int(i)
}

// apiTxn wraps the error handling of multiple transactions with the API; the
// user just checks the 'err' field at the end of the transaction block.
type apiTxn struct {
	err  error
	apis API
}

func (a *apiTxn) addCompany(co types.Company) (c *types.Company) {
	if a.err != nil {
		return &co
	}
	c, a.err = a.apis.Company.NewCompany(co)
	return c
}

func (a *apiTxn) getCompany(id int) (c *types.Company) {
	if a.err != nil {
		return c
	}
	var co types.Company
	co, a.err = a.apis.Company.GetCompany(id)
	return &co
}

func (a *apiTxn) addItem(item types.Item) (i *types.Item) {
	if a.err != nil {
		return &item
	}
	i, a.err = a.apis.Item.AddItem(item)
	return i
}

func (a *apiTxn) getItem(id int) (i *types.Item) {
	if a.err != nil {
		return i
	}
	var item types.Item
	item, a.err = a.apis.Item.GetItem(id)
	return &item

}

func (a *apiTxn) addPledge(itemID int, userID string) (p int) {
	if a.err != nil {
		return 0
	}
	p, a.err = a.apis.Pledge.NewPledge(itemID, userID)
	return p
}

func (a *apiTxn) deleteNewItem(id int) {
	if a.err != nil {
		return
	}
	a.err = apis.NewItem.DeleteNewItem(id)
}

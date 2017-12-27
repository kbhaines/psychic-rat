package pub

import (
	"log"
	"net/http"
	"psychic-rat/sess"
	"psychic-rat/types"
	"psychic-rat/web/dispatch"
	"strconv"

	"github.com/gorilla/sessions"
)

type (
	APIS struct {
		Item    ItemAPI
		NewItem NewItemAPI
		Pledge  PledgeAPI
		User    UserAPI
	}

	// TODO: APIs need consistent parameter and return style

	ItemAPI interface {
		ListItems() ([]types.Item, error)
		GetItem(id int) (types.Item, error)
	}

	NewItemAPI interface {
		AddNewItem(types.NewItem) (*types.NewItem, error)
	}

	PledgeAPI interface {
		AddPledge(itemId int, userId string) (int, error)
	}

	UserAPI interface {
		GetUser(userId string) (*types.User, error)
	}

	AuthHandler interface {
		Handler(http.ResponseWriter, *http.Request)
		GetLoggedInUser(*http.Request) (*types.User, error)
	}

	Renderer interface {
		Render(http.ResponseWriter, string, interface{}) error
	}

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
	apis APIS

	auth0Mode bool

	renderer Renderer

	// TODO: Env var
	auth0Store = sessions.NewCookieStore([]byte("something-very-secret"))

	authHandler AuthHandler
)

func Init(a APIS, useAuth0 bool, ah AuthHandler, rend Renderer) {
	apis = a
	auth0Mode = useAuth0
	authHandler = ah
	renderer = rend
}

func HomePageHandler(writer http.ResponseWriter, request *http.Request) {
	selector := dispatch.MethodSelector{
		"GET": func(writer http.ResponseWriter, request *http.Request) {
			vars := (&pageVariables{}).withSessionVars(request)
			renderer.Render(writer, "home.html.tmpl", vars)
		},
	}
	dispatch.ExecHandlerForMethod(selector, writer, request)
}

func SignInPageHandler(writer http.ResponseWriter, request *http.Request) {
	dispatch.ExecHandlerForMethod(dispatch.MethodSelector{"GET": authHandler.Handler}, writer, request)
}

func userLoginRequired(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := authHandler.GetLoggedInUser(r)
		if user == nil || err != nil {
			log.Print("user not logged in, or error occurred: %v", err)
			http.Error(w, "", http.StatusForbidden)
			return
		}
		h(w, r)
	}
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
	renderer.Render(writer, "pledge.html.tmpl", vars)
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
	renderer.Render(writer, "pledge-post.html.tmpl", vars)
}

func ThanksPageHandler(writer http.ResponseWriter, request *http.Request) {
	selector := dispatch.MethodSelector{
		"GET": userLoginRequired(func(writer http.ResponseWriter, request *http.Request) {
			vars := (&pageVariables{}).withSessionVars(request)
			renderer.Render(writer, "thanks.html.tmpl", vars)
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
	renderer.Render(w, "pledge-post-new-item.html.tmpl", vars)
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

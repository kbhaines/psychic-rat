package pub

import (
	"log"
	"net/http"
	"psychic-rat/types"
	"psychic-rat/web/dispatch"
	"strconv"
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
		AddPledge(itemId int, userId string) (*types.Pledge, error)
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

	// pageVariables holds data for the templates. Stuffed into one struct for now.
	pageVariables struct {
		Items     []types.Item
		User      types.User
		NewItems  []types.NewItem
		Companies []types.Company
	}
)

var (
	apis        APIS
	renderer    Renderer
	authHandler AuthHandler
)

func Init(a APIS, ah AuthHandler, rend Renderer) {
	apis = a
	authHandler = ah
	renderer = rend
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	selector := dispatch.MethodSelector{
		"GET": func(w http.ResponseWriter, r *http.Request) {
			vars := (&pageVariables{}).withSessionVars(r)
			renderer.Render(w, "home.html.tmpl", vars)
		},
	}
	dispatch.ExecHandlerForMethod(selector, w, r)
}

func SignInPageHandler(w http.ResponseWriter, r *http.Request) {
	dispatch.ExecHandlerForMethod(dispatch.MethodSelector{"GET": authHandler.Handler}, w, r)
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

func PledgePageHandler(w http.ResponseWriter, r *http.Request) {
	selector := dispatch.MethodSelector{
		"GET":  userLoginRequired(pledgeGetHandler),
		"POST": userLoginRequired(pledgePostHandler),
	}
	dispatch.ExecHandlerForMethod(selector, w, r)
}

func pledgeGetHandler(w http.ResponseWriter, r *http.Request) {
	report, err := apis.Item.ListItems()
	if err != nil {
		log.Print(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	vars := &pageVariables{Items: report}
	vars = vars.withSessionVars(r)
	renderer.Render(w, "pledge.html.tmpl", vars)
}

func pledgePostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Print(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	itemId64, err := strconv.ParseInt(r.FormValue("item"), 10, 32)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	itemId := int(itemId64)

	item, err := apis.Item.GetItem(itemId)
	if err != nil {
		log.Printf("error looking up item %v : %v", itemId, err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	user, _ := authHandler.GetLoggedInUser(r)
	userId := user.ID

	log.Printf("pledge item %v from user %v", itemId, userId)
	pledge, err := apis.Pledge.AddPledge(itemId, userId)
	if err != nil {
		log.Print("unable to pledge : ", err)
		return
	}
	log.Printf("pledge %v created", pledge.PledgeID)

	vars := &pageVariables{
		User:  *user,
		Items: []types.Item{item},
	}
	renderer.Render(w, "pledge-post.html.tmpl", vars)
}

func ThanksPageHandler(w http.ResponseWriter, r *http.Request) {
	selector := dispatch.MethodSelector{
		"GET": userLoginRequired(func(w http.ResponseWriter, r *http.Request) {
			vars := (&pageVariables{}).withSessionVars(r)
			renderer.Render(w, "thanks.html.tmpl", vars)
		}),
	}
	dispatch.ExecHandlerForMethod(selector, w, r)
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
	user, _ := authHandler.GetLoggedInUser(r)
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
	user, err := authHandler.GetLoggedInUser(r)
	if err != nil {
		// todo - return error?
		log.Fatal(err)
	}
	if user != nil {
		pv.User = *user
	}
	return pv
}

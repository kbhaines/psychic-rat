package pub

import (
	"net/http"
	"net/url"
	"psychic-rat/log"
	"psychic-rat/types"
	"psychic-rat/web/dispatch"
	"strconv"
)

type (
	// TODO: APIs need consistent parameter and return style
	ItemAPI interface {
		ListItems() ([]types.Item, error)
		GetItem(id int) (types.Item, error)
		ListCurrencies() ([]types.Currency, error)
	}

	NewItemAPI interface {
		AddNewItem(types.NewItem) (*types.NewItem, error)
	}

	PledgeAPI interface {
		AddPledge(itemId int, userID string, usdValue int) (*types.Pledge, error)
		ListUserPledges(userID string) ([]types.Pledge, error)
		ListRecentPledges(limit int) ([]types.RecentPledge, error)
		TotalPledges() (int, error)
	}

	UserAPI interface {
		GetUser(userID string) (*types.User, error)
	}

	UserHandler interface {
		GetLoggedInUser(*http.Request) (*types.User, error)
		GetUserCSRF(http.ResponseWriter, *http.Request) (string, error)
		LogOut(http.ResponseWriter, *http.Request) error
	}

	Renderer interface {
		Render(http.ResponseWriter, string, interface{}) error
	}

	HumanTester interface {
		IsHuman(f url.Values) error
	}

	// pageVariables holds data for the templates. Stuffed into one struct for now.
	pageVariables struct {
		Items         []types.Item
		User          types.User
		NewItems      []types.NewItem
		Companies     []types.Company
		Currencies    []types.Currency
		CSRFToken     string
		UserPledges   []types.Pledge
		RecentPledges []types.RecentPledge
		TotalPledges  int
	}
)

var (
	itemsAPI    ItemAPI
	newItemsAPI NewItemAPI
	pledgeAPI   PledgeAPI
	renderer    Renderer
	authHandler UserHandler
	humanTester HumanTester
)

func Init(item ItemAPI, newItems NewItemAPI, pledge PledgeAPI, auth UserHandler, rendr Renderer, ht HumanTester) {
	itemsAPI = item
	newItemsAPI = newItems
	pledgeAPI = pledge
	authHandler = auth
	renderer = rendr
	humanTester = ht
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" && r.RequestURI != "index.html" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	selector := dispatch.MethodSelector{
		"GET": func(w http.ResponseWriter, r *http.Request) {
			vars := (&pageVariables{}).withSessionVars(r)
			if vars.User.ID != "" {
				pledges, err := pledgeAPI.ListUserPledges(vars.User.ID)
				if err != nil {
					log.Errorf(r.Context(), "unable to retrieve pledges: %v", err)
					http.Error(w, "pledge service error", http.StatusInternalServerError)
					return
				}
				vars.UserPledges = pledges[0:3]
			}
			recentPledges, err := pledgeAPI.ListRecentPledges(10)
			if err != nil {
				log.Errorf(r.Context(), "unable to retrieve pledges: %v", err)
				http.Error(w, "pledge service error", http.StatusInternalServerError)
				return
			}
			vars.RecentPledges = recentPledges
			total, err := pledgeAPI.TotalPledges()
			if err != nil {
				log.Errorf(r.Context(), "unable to retrieve total: %v", err)
				http.Error(w, "pledge service error", http.StatusInternalServerError)
				return
			}
			vars.TotalPledges = total
			render(r, w, "home.html.tmpl", vars)
		},
	}
	dispatch.ExecHandlerForMethod(selector, w, r)
}

func SignInPageHandler(w http.ResponseWriter, r *http.Request) {
	selector := dispatch.MethodSelector{
		"GET": func(w http.ResponseWriter, r *http.Request) {
			vars := (&pageVariables{}).withSessionVars(r)
			render(r, w, "signin.html.tmpl", vars)
		},
	}
	dispatch.ExecHandlerForMethod(selector, w, r)
}

func SignOutPageHandler(w http.ResponseWriter, r *http.Request) {
	authHandler.LogOut(w, r)
}

func PledgePageHandler(w http.ResponseWriter, r *http.Request) {
	selector := dispatch.MethodSelector{
		"GET":  userLoginRequired(pledgeGetHandler),
		"POST": userLoginRequired(pledgePostHandler),
	}
	dispatch.ExecHandlerForMethod(selector, w, r)
}

func userLoginRequired(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := authHandler.GetLoggedInUser(r)
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}
		if err != nil {
			log.Errorf(r.Context(), "user (%v) not logged in, or error occurred: %v", user, err)
			http.Error(w, "", http.StatusForbidden)
			return
		}
		h(w, r)
	}
}

func pledgeGetHandler(w http.ResponseWriter, r *http.Request) {
	items, err := itemsAPI.ListItems()
	if err != nil {
		log.Errorf(r.Context(), "%v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	currencies, err := itemsAPI.ListCurrencies()
	if err != nil {
		log.Errorf(r.Context(), "%v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	token, err := authHandler.GetUserCSRF(w, r)
	if err != nil {
		log.Errorf(r.Context(), "unable to get CSRF for user: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	vars := &pageVariables{Items: items, Currencies: currencies, CSRFToken: token}
	vars = vars.withSessionVars(r)
	render(r, w, "pledge.html.tmpl", vars)
}

func pledgePostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Errorf(r.Context(), "could not parse pledge form: %v", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := humanTester.IsHuman(r.Form); err != nil {
		log.Errorf(r.Context(), "failed human test: %v", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	itemId64, err := strconv.ParseInt(r.FormValue("item"), 10, 32)
	if err != nil {
		log.Errorf(r.Context(), "could not parse itemID: %v", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	itemId := int(itemId64)

	item, err := itemsAPI.GetItem(itemId)
	if err != nil {
		log.Errorf(r.Context(), "could not get item %v: %v", itemId, err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	user, _ := authHandler.GetLoggedInUser(r)
	userID := user.ID

	log.Logf(r.Context(), "pledge item %v from user %v", itemId, userID)

	pledges, err := pledgeAPI.ListUserPledges(userID)
	if err != nil {
		log.Errorf(r.Context(), "could not get pledges for user %v: %v", userID, err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if len(pledges) > 999999 {
		log.Errorf(r.Context(), "user pledge limit exceeded for %s", userID)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	pledge, err := pledgeAPI.AddPledge(itemId, userID, item.USDValue)
	if err != nil {
		log.Errorf(r.Context(), "unable to pledge: ", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	log.Logf(r.Context(), "pledge %v created", pledge.PledgeID)

	http.Redirect(w, r, "/thanks", http.StatusSeeOther)
}

func ThanksPageHandler(w http.ResponseWriter, r *http.Request) {
	selector := dispatch.MethodSelector{
		"GET": userLoginRequired(func(w http.ResponseWriter, r *http.Request) {
			vars := (&pageVariables{}).withSessionVars(r)
			render(r, w, "thanks.html.tmpl", vars)
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
		log.Errorf(r.Context(), "could not parse new item form: %v", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := humanTester.IsHuman(r.Form); err != nil {
		log.Errorf(r.Context(), "failed human test: %v", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	user, _ := authHandler.GetLoggedInUser(r)
	userID := user.ID

	company := r.FormValue("company")
	make := r.FormValue("make")
	model := r.FormValue("model")
	currencyID := r.FormValue("currencyID")
	value := r.FormValue("value")

	if company == "" || model == "" || make == "" || currencyID == "" || value == "" {
		log.Errorf(r.Context(), "new item request missing data: %v", r.Form)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	valueInt, err := strconv.ParseInt(r.FormValue("value"), 10, 32)
	if err != nil {
		log.Errorf(r.Context(), "could not parse value of new item: %v", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	currencyIDInt, err := strconv.ParseInt(r.FormValue("currencyID"), 10, 32)
	if err != nil {
		log.Errorf(r.Context(), "could not parse value of currencyID: %v", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	newItem := types.NewItem{UserID: userID, IsPledge: true, Make: make, Model: model, Company: company, CurrencyID: int(currencyIDInt), Value: int(valueInt)}
	_, err = newItemsAPI.AddNewItem(newItem)
	if err != nil {
		log.Errorf(r.Context(), "could not add new item %v: %v", newItem, err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	item := types.Item{Make: newItem.Make, Model: newItem.Model, Company: types.Company{Name: company}}
	vars := &pageVariables{
		User:  *user,
		Items: []types.Item{item},
	}
	render(r, w, "pledge-post-new-item.html.tmpl", vars)
}

func (pv *pageVariables) withSessionVars(r *http.Request) *pageVariables {
	user, err := authHandler.GetLoggedInUser(r)
	if err != nil {
		log.Errorf(r.Context(), "error retrieving user: %v", err)
		pv.User = types.User{}
		return pv
	}
	if user != nil {
		pv.User = *user
	}
	return pv
}

func render(r *http.Request, w http.ResponseWriter, template string, vars *pageVariables) {
	err := renderer.Render(w, template, vars)
	if err != nil {
		log.Errorf(r.Context(), "could not render template %s: %v", template, err)
		http.Error(w, "", http.StatusInternalServerError)
	}
	return
}

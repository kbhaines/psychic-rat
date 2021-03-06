package admin

import (
	"net/http"
	"psychic-rat/log"
	"psychic-rat/types"
	"psychic-rat/web/dispatch"
)

type (
	newItemPostData struct {
		ID          int
		Add         bool
		Delete      bool
		Pledge      bool
		ItemID      int
		CompanyID   int
		UserID      string
		UserCompany string
		UserMake    string
		UserModel   string
		Value       int
		CurrencyID  int
	}

	CompanyAPI interface {
		ListCompanies() ([]types.Company, error)
		GetCompany(int) (types.Company, error)
		AddCompany(co types.Company) (*types.Company, error)
	}

	ItemAPI interface {
		ListItems() ([]types.Item, error)
		AddItem(types.Item) (*types.Item, error)
		GetItem(id int) (types.Item, error)
		CurrencyConversion(id int, value int) (int, error)
	}

	NewItemAPI interface {
		ListNewItems() ([]types.NewItem, error)
		DeleteNewItem(int) error
		MarkNewItemUsed(int) error
	}

	PledgeAPI interface {
		AddPledge(itemId int, userId string, usdValue int) (*types.Pledge, error)
	}

	Renderer interface {
		Render(w http.ResponseWriter, templateName string, variables interface{}) error
	}

	UserHandler interface {
		GetLoggedInUser(*http.Request) (*types.User, error)
		GetUserCSRF(http.ResponseWriter, *http.Request) (string, error)
	}

	pageVariables struct {
		Items     []types.Item
		NewItems  []types.NewItem
		Companies []types.Company
		User      types.User
		CSRFToken string
	}
)

var (
	companyAPI  CompanyAPI
	itemsAPI    ItemAPI
	newItemsAPI NewItemAPI
	pledgeAPI   PledgeAPI
	renderer    Renderer
	authHandler UserHandler
)

func Init(co CompanyAPI, item ItemAPI, newItems NewItemAPI, pledge PledgeAPI, auth UserHandler, rendr Renderer) {
	companyAPI = co
	itemsAPI = item
	newItemsAPI = newItems
	pledgeAPI = pledge
	renderer = rendr
	authHandler = auth
}

func AdminItemHandler(w http.ResponseWriter, r *http.Request) {
	selector := dispatch.MethodSelector{
		"GET":  adminLoginRequired(listNewItems),
		"POST": adminLoginRequired(approveNewItems),
	}
	dispatch.ExecHandlerForMethod(selector, w, r)
}

func adminLoginRequired(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isUserAdmin(r) {
			log.Errorf(r.Context(), "user is not admin")
			http.Error(w, "", http.StatusForbidden)
			return
		}
		h(w, r)
	}
}

func isUserAdmin(r *http.Request) bool {
	user, err := authHandler.GetLoggedInUser(r)
	if err != nil {
		log.Errorf(r.Context(), "could not retrieve user: %v", err)
		return false
	}
	return user != nil && user.IsAdmin
}

func listNewItems(w http.ResponseWriter, r *http.Request) {
	newItems, err := newItemsAPI.ListNewItems()
	if err != nil {
		log.Errorf(r.Context(), "could not retrieve new items: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	items, err := itemsAPI.ListItems()
	if err != nil {
		log.Errorf(r.Context(), "could not retrieve items: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	companies, err := companyAPI.ListCompanies()
	if err != nil {
		log.Errorf(r.Context(), "could not retrieve companies: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	token, err := authHandler.GetUserCSRF(w, r)
	if err != nil {
		log.Errorf(r.Context(), "unable to get CSRF for user: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	vars := &pageVariables{Items: items, NewItems: newItems, Companies: companies, CSRFToken: token}
	renderer.Render(w, "admin-new-items.html.tmpl", vars)
}

func approveNewItems(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Errorf(r.Context(), "could not parse form: %v", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	reader := newFormReader(r.Form)
	for reader.next() {
		nip := reader.getNewItemPost()
		if reader.errors() {
			break
		}
		if err := processNewItemPost(r, nip); err != nil {
			log.Errorf(r.Context(), "could not complete transactions for new item %d:  %v", nip.ID, err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	if reader.errors() {
		log.Errorf(r.Context(), "errors while parsing form line %d: %v", reader.row, reader.err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}

func processNewItemPost(r *http.Request, nip newItemPostData) error {
	if nip.Delete {
		log.Logf(r.Context(), "deleting new item %v", nip)
		return newItemsAPI.DeleteNewItem(nip.ID)
	}
	if !nip.Add {
		log.Logf(r.Context(), "ignoring new item %v", nip)
		return nil
	}

	txn := apiTxn{nil}
	var item *types.Item
	var value int
	if nip.ItemID == 0 {
		log.Logf(r.Context(), "creating new item from %v", nip)
		var company *types.Company
		if nip.CompanyID == 0 {
			log.Logf(r.Context(), "creating new company from %v", nip.UserCompany)
			company = txn.addCompany(types.Company{Name: nip.UserCompany})
		} else {
			log.Logf(r.Context(), "using existing company %v", nip.CompanyID)
			company = txn.getCompany(nip.CompanyID)
		}
		value = txn.currencyConversion(nip.CurrencyID, nip.Value)
		item = txn.addItem(types.Item{Company: *company, Make: nip.UserMake, Model: nip.UserModel, USDValue: value, NewItemID: nip.ID})
		log.Logf(r.Context(), "item added: %v", item)
	} else {
		log.Logf(r.Context(), "using existing item %d", nip.ItemID)
		item = txn.getItem(nip.ItemID)
	}
	txn.markUsed(nip.ID)
	if nip.Pledge {
		txn.addPledge(item, nip.UserID, value)
	}
	return txn.err
}

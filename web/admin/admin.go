package admin

import (
	"log"
	"net/http"
	"psychic-rat/types"
	"psychic-rat/web/dispatch"
)

type (
	NewItemAdminPost struct {
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
	}

	NewItemAPI interface {
		ListNewItems() ([]types.NewItem, error)
		DeleteNewItem(int) error
	}
	PledgeAPI interface {
		AddPledge(itemId int, userId string) (*types.Pledge, error)
	}

	Renderer interface {
		Render(w http.ResponseWriter, templateName string, variables interface{}) error
	}

	AuthHandler interface {
		GetLoggedInUser(*http.Request) (*types.User, error)
	}

	pageVariables struct {
		Items     []types.Item
		NewItems  []types.NewItem
		Companies []types.Company
		User      types.User
	}
)

var (
	companyAPI  CompanyAPI
	itemsAPI    ItemAPI
	newItemsAPI NewItemAPI
	pledgeAPI   PledgeAPI
	renderer    Renderer
	authHandler AuthHandler
)

func Init(co CompanyAPI, item ItemAPI, newItems NewItemAPI, pledge PledgeAPI, auth AuthHandler, rendr Renderer) {
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
			log.Print("user not admin")
			http.Error(w, "", http.StatusForbidden)
			return
		}
		h(w, r)
	}
}

func isUserAdmin(request *http.Request) bool {
	user, err := authHandler.GetLoggedInUser(request)
	if err != nil {
		log.Print(err)
		return false
	}
	return user != nil && user.IsAdmin
}

func listNewItems(w http.ResponseWriter, r *http.Request) {
	newItems, err := newItemsAPI.ListNewItems()
	if err != nil {
		log.Printf("unable to retrieve new items: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	items, err := itemsAPI.ListItems()
	if err != nil {
		log.Printf("unable to retrieve items: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	companies, err := companyAPI.ListCompanies()
	if err != nil {
		log.Printf("unable to retrieve companies: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	vars := &pageVariables{Items: items, NewItems: newItems, Companies: companies}
	renderer.Render(w, "admin-new-items.html.tmpl", vars)
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
		if err := processNewItemPost(nip); err != nil {
			log.Printf("unable to complete transactions for new item %d:  %v", nip.ID, err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

	}
	if reader.errors() {
		log.Printf("errors while parsing form line %d: %v", reader.row, reader.err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}

func processNewItemPost(nip NewItemAdminPost) error {
	if nip.Delete {
		return newItemsAPI.DeleteNewItem(nip.ID)
	}
	txn := apiTxn{nil}
	var item *types.Item
	if nip.ItemID == 0 {
		var company *types.Company
		if nip.CompanyID == 0 {
			company = txn.addCompany(types.Company{Name: nip.UserCompany})
		} else {
			company = txn.getCompany(nip.CompanyID)
		}
		item = txn.addItem(types.Item{Company: *company, Make: nip.UserMake, Model: nip.UserModel})
	} else {
		item = txn.getItem(nip.ItemID)
	}

	if nip.Pledge {
		txn.addPledge(item.ID, nip.UserID)
	}
	txn.deleteNewItem(nip.ID)
	return txn.err
}

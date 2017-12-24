package admin

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"psychic-rat/sess"
	"psychic-rat/types"
	"psychic-rat/web/dispatch"
	"psychic-rat/web/tmpl"
	"strconv"
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

	// apiTxn wraps the error handling of multiple transactions with the API; the
	// user just checks the 'err' field at the end of the transaction block.
	apiTxn struct {
		err error
	}

	// formReader parses a submitted New Items form POST request, captures multiple
	// errors that resulted from parsing.
	formReader struct {
		form url.Values
		row  int
		rows []int
		err  []error
	}

	CompanyAPI interface {
		GetCompanies() ([]types.Company, error)
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
		AddPledge(itemId int, userId string) (int, error)
	}

	pageVariables struct {
		Items     []types.Item
		NewItems  []types.NewItem
		Companies []types.Company
	}
)

var (
	companyAPI  CompanyAPI
	itemsAPI    ItemAPI
	newItemsAPI NewItemAPI
	pledgeAPI   PledgeAPI
	renderPage  = tmpl.RenderTemplate
)

func InitDeps(co CompanyAPI, item ItemAPI, newItems NewItemAPI, pledge PledgeAPI) {
	companyAPI = co
	itemsAPI = item
	newItemsAPI = newItems
	pledgeAPI = pledge
}

func AdminItemHandler(w http.ResponseWriter, r *http.Request) {
	selector := map[string]http.HandlerFunc{
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
	s := sess.NewSessionStore(request, nil)
	user, err := s.Get()
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
	companies, err := companyAPI.GetCompanies()
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

		if err := processNewItemPost(nip); err != nil {
			log.Printf("unable to complete transactions for new item %d:  %v", nip.ID, err)
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

func processNewItemPost(nip NewItemAdminPost) error {
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

func newFormReader(form url.Values) *formReader {
	fr := &formReader{form, -1, []int{}, []error{}}
	adds, ok := form["add[]"]
	if !ok {
		return fr
	}
	for _, str := range adds {
		rowID, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			log.Printf("unable to parse row Id: %v", err)
			fr.err = append(fr.err, err)
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
	if f.errors() {
		panic("getNewItemPost called when formReader in error state")
	}

	i := NewItemAdminPost{
		ID:          f.getInt("id[]"),
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

func (a *apiTxn) addCompany(co types.Company) (c *types.Company) {
	if a.err != nil {
		return &co
	}
	c, a.err = companyAPI.AddCompany(co)
	return c
}

func (a *apiTxn) getCompany(id int) (c *types.Company) {
	if a.err != nil {
		return c
	}
	var co types.Company
	co, a.err = companyAPI.GetCompany(id)
	return &co
}

func (a *apiTxn) addItem(item types.Item) (i *types.Item) {
	if a.err != nil {
		return &item
	}
	i, a.err = itemsAPI.AddItem(item)
	return i
}

func (a *apiTxn) getItem(id int) (i *types.Item) {
	if a.err != nil {
		return i
	}
	var item types.Item
	item, a.err = itemsAPI.GetItem(id)
	return &item

}

func (a *apiTxn) addPledge(itemID int, userID string) (p int) {
	if a.err != nil {
		return 0
	}
	p, a.err = pledgeAPI.AddPledge(itemID, userID)
	return p
}

func (a *apiTxn) deleteNewItem(id int) {
	if a.err != nil {
		return
	}
	a.err = newItemsAPI.DeleteNewItem(id)
}

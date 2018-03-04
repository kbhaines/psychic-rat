package tests

import (
	"fmt"
	"net/http"
	"net/url"
	"psychic-rat/sqldb"
	"psychic-rat/types"
	"reflect"
	"strconv"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

type (
	// newItemHtml holds a row's worth of data scraped from the admin/items
	// web page
	newItemHtml struct {
		ID          string
		IsPledge    string
		UserID      string
		UserCompany string
		UserMake    string
		UserModel   string
		UserValue   string
		CurrencyID  string
	}

	// postLine allows us to build up an array-based POST request
	postLine struct {
		row int        // row is the active row we're populating
		v   url.Values // v is the values we'll post eventually
	}

	newItemPost struct {
		Company    string
		Make       string
		Model      string
		CurrencyID string
		Value      string
	}
)

func TestBlockAccessToItemListing(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	resp, err := getAuthdClient("test1", t).Get(testUrl + "/admin/newitems")
	testPageStatus(resp, err, http.StatusForbidden, t)
}

func TestListNewItems(t *testing.T) {
	server, db := newServer(t)
	initNewItems(db, t)
	defer cleanUp(server, db)
	testNewItemsPage(testNewItems, t)
}

func testNewItemsPage(expectedNewItems []types.NewItem, t *testing.T) {
	t.Helper()
	resp, err := getAuthdClient("admin", t).Get(testUrl + "/admin/newitems")
	testPageStatus(resp, err, http.StatusOK, t)

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		t.Fatal(err)
	}

	fmap := map[string]func(*newItemHtml, string){
		"id[]":          (*newItemHtml).id,
		"isPledge[]":    (*newItemHtml).pledge,
		"userID[]":      (*newItemHtml).userID,
		"usercompany[]": (*newItemHtml).userCompany,
		"usermake[]":    (*newItemHtml).userMake,
		"usermodel[]":   (*newItemHtml).userModel,
		"uservalue[]":   (*newItemHtml).userValue,
		"action[]":      (*newItemHtml).nowt,
		"currencyID[]":  (*newItemHtml).currencyID,
	}

	rows := doc.Find(".items-table .item-entry")
	if rows.Size() != len(expectedNewItems) {
		t.Fatalf("expected %d rows in item listing, got %d", len(expectedNewItems), rows.Size())
	}
	rows.Each(func(i int, s *goquery.Selection) {
		actualNewItem := newItemHtml{}

		s.Find("input").Each(func(_ int, s *goquery.Selection) {
			n := s.AttrOr("name", "")
			v := s.AttrOr("value", "")
			f, ok := fmap[n]
			if !ok {
				t.Fatalf("unknown input: %s", n)
			}
			f(&actualNewItem, v)
		})

		expectedNewItem := newItemHtml{
			ID:          strconv.Itoa(i + 1),
			IsPledge:    "true",
			UserID:      expectedNewItems[i].UserID,
			UserCompany: expectedNewItems[i].Company,
			UserMake:    expectedNewItems[i].Make,
			UserModel:   expectedNewItems[i].Model,
			UserValue:   "100",
			CurrencyID:  strconv.Itoa(expectedNewItems[i].CurrencyID),
		}
		if !reflect.DeepEqual(expectedNewItem, actualNewItem) {
			t.Errorf("expected html form to have %v, got %v", expectedNewItem, actualNewItem)
		}
	})

}

func TestBadnewItemPost(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	client := getAuthdClient("admin", t)
	post := url.Values{
		"id[]":       {"0"},
		"action[]":   {"blah"},
		"isPledge[]": {"0"},
		"userID[]":   {"test1"},
		"csrf":       {getCSRFToken(client, testUrl+"/admin/newitems", t)},
	}
	resp, err := client.PostForm(testUrl+"/admin/newitems", post)
	testPageStatus(resp, err, http.StatusBadRequest, t)
}

func TestNewItemAdminPost(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	initNewItems(db, t)
	// Going backwards so we don't have to adjust the index for each row
	// and can just use the coincident index values from the DB. Bit lazy....
	client := getAuthdClient("admin", t)
	for i := len(testNewItems) - 1; i > 0; i-- {
		ni := testNewItems[i]
		pl := &postLine{v: url.Values{}}
		pl = pl.newPostLine(ni.ID).
			userID(ni.UserID).
			userCompany(ni.Company).
			userMake(ni.Make).
			userModel(ni.Model).
			currency(ni.CurrencyID).
			value(ni.Value).
			isPledge().
			selectToAdd()

		pl.v.Add("csrf", getCSRFToken(client, testUrl+"/admin/newitems", t))
		resp, err := client.PostForm(testUrl+"/admin/newitems", pl.v)
		testPageStatus(resp, err, http.StatusOK, t)

		item := findAddedNewItem(ni, db, t)
		testNewItemPledged(ni.UserID, *item, db, t)
		testNewItemsPage(testNewItems[:i], t)
	}
}

func findAddedNewItem(item types.NewItem, db *sqldb.DB, t *testing.T) *types.Item {
	t.Helper()
	items, err := db.ListItems()
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range items {
		if i.Make == item.Make && i.Model == item.Model && i.Company.Name == item.Company && i.USDValue == item.Value {
			return &i
		}
	}
	t.Fatalf("did not find expected item %v in items db:\n%v", item, items)
	return nil
}

func testNewItemPledged(userID string, item types.Item, db *sqldb.DB, t *testing.T) {
	t.Helper()
	pledges, err := db.ListUserPledges(userID)
	if err != nil {
		t.Fatalf("could not list pledges for %s: %v", "test1", err)
		return
	}
	for _, p := range pledges {
		pi := p.Item
		if pi.Make == item.Make && pi.Model == item.Model && pi.Company.ID == item.Company.ID {
			return
		}
	}
	t.Fatalf("expected to find %v in pledges for user %s, found %v", item, "test1", pledges)
}

func TestNewItemAdminPostUsingExistingItem(t *testing.T) {
	server, db := newServer(t)
	initNewItems(db, t)
	defer cleanUp(server, db)
	currentItems, err := db.ListItems()
	if err != nil {
		t.Fatal(err)
	}

	pl := postLine{v: url.Values{}}
	pl.newPostLine(6).userID(testNewItems[5].UserID).existingItem(1).selectToAdd()

	client := getAuthdClient("admin", t)
	pl.v.Add("csrf", getCSRFToken(client, testUrl+"/admin/newitems", t))
	resp, err := client.PostForm(testUrl+"/admin/newitems", pl.v)

	newItems, err := db.ListItems()
	if len(currentItems) != len(newItems) {
		t.Fatalf("expected %d items, got %d items", len(currentItems), len(newItems))
	}

	testPageStatus(resp, err, http.StatusOK, t)
	testNewItemsPage(testNewItems[:len(testNewItems)-1], t)
}

func TestNewItemAdminPostUsingExistingCompany(t *testing.T) {
	server, db := newServer(t)
	initNewItems(db, t)
	defer cleanUp(server, db)
	currentCompanies, err := db.ListCompanies()
	if err != nil {
		t.Fatal(err)
	}

	pl := postLine{v: url.Values{}}
	ni := testNewItems[5]

	pl.newPostLine(6).userID(ni.UserID).existingCompany(1).currency(ni.CurrencyID).value(ni.Value).selectToAdd()
	client := getAuthdClient("admin", t)
	pl.v.Add("csrf", getCSRFToken(client, testUrl+"/admin/newitems", t))
	resp, err := client.PostForm(testUrl+"/admin/newitems", pl.v)

	newCompanies, err := db.ListCompanies()
	if len(currentCompanies) != len(newCompanies) {
		t.Fatalf("expected %d items, got %d items", len(currentCompanies), len(newCompanies))
	}

	testPageStatus(resp, err, http.StatusOK, t)
	testNewItemsPage(testNewItems[:len(testNewItems)-1], t)
}

func TestDeleteNewItems(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	initNewItems(db, t)
	pl := postLine{v: url.Values{}}
	ni := testNewItems[5]
	pl.newPostLine(6).userID(ni.UserID).selectToDelete()
	client := getAuthdClient("admin", t)
	pl.v.Add("csrf", getCSRFToken(client, testUrl+"/admin/newitems", t))
	resp, err := client.PostForm(testUrl+"/admin/newitems", pl.v)

	testPageStatus(resp, err, http.StatusOK, t)
	testNewItemsPage(testNewItems[:len(testNewItems)-1], t)
}

func TestBadNewItemsPostInvalidAddParams(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	v := url.Values{"action[]": []string{"abc"}, "id[]": []string{"abc"}}
	client := getAuthdClient("admin", t)
	v.Add("csrf", getCSRFToken(client, testUrl+"/admin/newitems", t))
	resp, err := client.PostForm(testUrl+"/admin/newitems", v)
	testPageStatus(resp, err, http.StatusBadRequest, t)
}

func (n *newItemHtml) id(w string)          { n.ID = w }
func (n *newItemHtml) pledge(w string)      { n.IsPledge = w }
func (n *newItemHtml) userID(w string)      { n.UserID = w }
func (n *newItemHtml) userCompany(w string) { n.UserCompany = w }
func (n *newItemHtml) userMake(w string)    { n.UserMake = w }
func (n *newItemHtml) userModel(w string)   { n.UserModel = w }
func (n *newItemHtml) userValue(w string)   { n.UserValue = w }
func (n *newItemHtml) currencyID(w string)  { n.CurrencyID = w }
func (n *newItemHtml) nowt(w string)        {}

func spfi(i int) string { return fmt.Sprintf("%d", i) }

func (p *postLine) newPostLine(itemID int) *postLine {
	p.row = len(p.v["id[]"])
	p.v.Add("id[]", spfi(itemID))
	p.v.Add("action[]", "leave")
	p.v.Add("isPledge[]", "0")
	p.v.Add("item[]", "0")
	p.v.Add("company[]", "0")
	p.v.Add("userID[]", "")
	p.v.Add("usercompany[]", "")
	p.v.Add("usermake[]", "")
	p.v.Add("usermodel[]", "")
	p.v.Add("uservalue[]", "0")
	p.v.Add("currencyID[]", "0")
	return p
}

func (p *postLine) selectToAdd() *postLine          { p.v["action[]"][p.row] = "add"; return p }
func (p *postLine) selectToDelete() *postLine       { p.v["action[]"][p.row] = "delete"; return p }
func (p *postLine) isPledge() *postLine             { p.v["isPledge[]"][p.row] = "true"; return p }
func (p *postLine) userID(u string) *postLine       { p.v["userID[]"][p.row] = u; return p }
func (p *postLine) existingItem(i int) *postLine    { p.v["item[]"][p.row] = spfi(i); return p }
func (p *postLine) existingCompany(c int) *postLine { p.v["company[]"][p.row] = spfi(c); return p }
func (p *postLine) userCompany(c string) *postLine  { p.v["usercompany[]"][p.row] = c; return p }
func (p *postLine) userMake(m string) *postLine     { p.v["usermake[]"][p.row] = m; return p }
func (p *postLine) userModel(m string) *postLine    { p.v["usermodel[]"][p.row] = m; return p }
func (p *postLine) value(v int) *postLine           { p.v["uservalue[]"][p.row] = spfi(v); return p }
func (p *postLine) currency(c int) *postLine        { p.v["currencyID[]"][p.row] = spfi(c); return p }

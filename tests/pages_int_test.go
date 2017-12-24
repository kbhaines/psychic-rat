package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"psychic-rat/sqldb"
	"psychic-rat/web"
	"psychic-rat/web/admin"
	"psychic-rat/web/tmpl"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

// newItemHtml holds a row's worth of data scraped from the admin/items
// web page
type (
	newItemHtml struct {
		ID          string
		IsPledge    string
		UserID      string
		UserCompany string
		UserMake    string
		UserModel   string
		UserValue   string
	}

	// postLine allows us to build up an array-based POST request
	postLine struct {
		row int        // row is the active row we're populating
		v   url.Values // v is the values we'll post eventually
	}

	newItemPost struct {
		Company string
		Make    string
		Model   string
	}
)

var (
	testUrl string

	newItemPosts = []newItemPost{
		newItemPost{Company: "newco1", Make: "newmake1", Model: "newmodel1"},
		newItemPost{Company: "newco2", Make: "newmake2", Model: "newmodel2"},
		newItemPost{Company: "newco3", Make: "newmake3", Model: "newmodel3"},
		newItemPost{Company: "newco4", Make: "newmake4", Model: "newmodel4"},
	}
)

func TestHomePage(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	resp, err := http.Get(testUrl + "/")
	testPageStatus(resp, err, http.StatusOK, t)
}

func newServer(t *testing.T) (*httptest.Server, *sqldb.DB) {
	server := httptest.NewServer(web.Handler())
	testUrl = server.URL
	db := initDB(t)
	apis := web.APIS{
		Company: db,
		Item:    db,
		NewItem: db,
		Pledge:  db,
		User:    db,
	}
	web.Init(apis)
	admin.Init(db, db, db, db)
	tmpl.Init("../res/")
	return server, db
}

func testPageStatus(resp *http.Response, err error, expectedCode int, t *testing.T) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error", err)
	}
	if resp.StatusCode != expectedCode {
		t.Fatalf("wanted %v, got %v. Response was %v", expectedCode, resp.StatusCode, resp)
	}
}

func cleanUp(server *httptest.Server, db *sqldb.DB) {
	server.Close()
	closeDB(db)
}

func TestPledgeWithoutLogin(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	resp, err := http.Get(testUrl + "/pledge")
	testPageStatus(resp, err, http.StatusForbidden, t)
}

func TestThankYouWithoutLogin(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	resp, err := http.Get(testUrl + "/thanks")
	testPageStatus(resp, err, http.StatusForbidden, t)
}

func TestSignin(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	loginUser("test1", t)
}

func loginUser(user string, t *testing.T) http.CookieJar {
	t.Helper()
	resp, err := http.Get(testUrl + "/signin?u=" + user)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("unable to signin, error was %v, response %v", err, resp)
	}
	jar, _ := cookiejar.New(nil)
	url, _ := url.Parse(testUrl)
	jar.SetCookies(url, resp.Cookies())
	return jar
}

func TestPledgeWithLogin(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)

	cookie := loginUser("test1", t)
	client := http.Client{Jar: cookie}
	resp, err := client.Get(testUrl + "/pledge")
	testPageStatus(resp, err, http.StatusOK, t)

	strBody := readResponseBody(resp, t)
	expected := []string{
		"user1 full",
		"<select ",
		"<input type=\"submit\"",
	}
	testStrings(strBody, expected, t)
}

func readResponseBody(resp *http.Response, t *testing.T) string {
	t.Helper()
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

func testStrings(body string, expectedStrings []string, t *testing.T) {
	t.Helper()
	for _, s := range expectedStrings {
		if !strings.Contains(body, s) {
			t.Errorf("body did not contain '%s'", s)
		}
	}
	if t.Failed() {
		t.Errorf("body was %s", body)
	}
}

func TestHappyPathPledge(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	cookie := loginUser("test1", t)
	client := http.Client{Jar: cookie}
	data := url.Values{"item": {"1"}}
	resp, err := client.PostForm(testUrl+"/pledge", data)
	testPageStatus(resp, err, http.StatusOK, t)

	expected := []string{
		"boycott of phone xyz by testco1",
		"Signed in as user1 full",
	}
	body := readResponseBody(resp, t)
	testStrings(body, expected, t)
}

func TestBadNewItems(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	cookie := loginUser("test1", t)
	client := http.Client{Jar: cookie}

	values := []url.Values{
		url.Values{"company": {""}, "make": {"newmake"}, "model": {"newmodel"}},
		url.Values{"company": {"newco"}, "make": {""}, "model": {"newmodel"}},
		url.Values{"company": {""}, "make": {"newmake"}, "model": {""}},
		url.Values{"make": {"newmake"}, "model": {"bla"}},
	}

	for _, d := range values {
		resp, err := client.PostForm(testUrl+"/newitem", d)
		testPageStatus(resp, err, http.StatusBadRequest, t)
	}
	items, err := db.ListNewItems()
	if err != nil {
		t.Fatal(err)
	}
	if 0 != len(items) {
		t.Fatalf("expected 0 new items, got %d", len(items))
	}
}

func TestBlockAccessToItemListing(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	cookie := loginUser("test1", t)
	client := http.Client{Jar: cookie}

	resp, err := client.Get(testUrl + "/admin/newitems")
	testPageStatus(resp, err, http.StatusForbidden, t)
}

func TestListNewItems(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	cookie := loginUser("test1", t)
	client := http.Client{Jar: cookie}
	loadNewItems(client, t)

	testNewItemsPage(newItemPosts, t)
}

func testNewItemsPage(newItems []newItemPost, t *testing.T) {
	t.Helper()
	cookie := loginUser("admin", t)
	client := http.Client{Jar: cookie}
	resp, err := client.Get(testUrl + "/admin/newitems")
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
		"add[]":         (*newItemHtml).nowt,
		"delete[]":      (*newItemHtml).nowt,
	}

	rows := doc.Find(".items-table .item-entry")
	if rows.Size() != len(newItems) {
		t.Fatalf("expected %d rows in item listing, got %d", len(newItems), rows.Size())
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
			UserID:      "test1",
			UserCompany: newItems[i].Company,
			UserMake:    newItems[i].Make,
			UserModel:   newItems[i].Model,
			UserValue:   "0",
		}
		if !reflect.DeepEqual(expectedNewItem, actualNewItem) {
			t.Errorf("expected html form to have %v, got %v", expectedNewItem, actualNewItem)
		}
	})

}

func loadNewItems(client http.Client, t *testing.T) {
	t.Helper()
	for i, d := range newItemPosts {
		resp, err := client.PostForm(testUrl+"/newitem", d.getUrlValues())
		testPageStatus(resp, err, http.StatusOK, t)
		expected := []string{
			fmt.Sprintf("boycott of %s %s by %s", newItemPosts[i].Make, newItemPosts[i].Model, newItemPosts[i].Company),
			"Signed in as user1 full",
			"new item is under review",
		}
		body := readResponseBody(resp, t)
		testStrings(body, expected, t)
	}
}

func (n *newItemPost) getUrlValues() url.Values {
	return url.Values{"company": {n.Company}, "make": {n.Make}, "model": {n.Model}}
}

func TestBadnewItemPost(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	cookie := loginUser("admin", t)
	client := http.Client{Jar: cookie}

	post := url.Values{
		"add[]":      []string{"0"}, // row selected
		"isPledge[]": []string{"0"},
		"userID[]":   []string{"test1"}, // user test1
	}
	resp, err := client.PostForm(testUrl+"/admin/newitems", post)
	testPageStatus(resp, err, http.StatusBadRequest, t)

}

func TestNewItemAdminPost(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	cookie := loginUser("test1", t)
	client := http.Client{Jar: cookie}
	loadNewItems(client, t)

	cookie = loginUser("admin", t)
	client = http.Client{Jar: cookie}

	// Going backwards so we don't have to adjust the index for each row
	// and can just use the coincident index values from the DB. Bit lazy....
	for itemToUse := 3; itemToUse > 0; itemToUse-- {
		pl := &postLine{v: url.Values{}}
		pl = pl.newPostLine(itemToUse + 1).userID("test1").
			userCompany(newItemPosts[itemToUse].Company).
			userMake(newItemPosts[itemToUse].Make).
			userModel(newItemPosts[itemToUse].Model).
			selectToAdd()

		resp, err := client.PostForm(testUrl+"/admin/newitems", pl.v)
		testPageStatus(resp, err, http.StatusOK, t)

		testNewItemAdded(newItemPosts[itemToUse], db, t)
		testNewItemsPage(newItemPosts[:itemToUse], t)
	}
}

func testNewItemAdded(item newItemPost, db *sqldb.DB, t *testing.T) {
	t.Helper()
	items, err := db.ListItems()
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range items {
		if i.Make == item.Make && i.Model == item.Model && i.Company.Name == item.Company {
			return
		}
	}
	t.Fatalf("did not find expected item %v in items db", item)
}

func TestNewItemAdminPostUsingExistingItem(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	cookie := loginUser("test1", t)
	client := http.Client{Jar: cookie}
	loadNewItems(client, t)

	cookie = loginUser("admin", t)
	client = http.Client{Jar: cookie}

	currentItems, err := db.ListItems()
	if err != nil {
		t.Fatal(err)
	}

	pl := postLine{v: url.Values{}}
	pl.newPostLine(4).userID("test1").existingItem(1).selectToAdd()
	resp, err := client.PostForm(testUrl+"/admin/newitems", pl.v)

	newItems, err := db.ListItems()
	if len(currentItems) != len(newItems) {
		t.Fatalf("expected %d items, got %d items", len(currentItems), len(newItems))
	}

	testPageStatus(resp, err, http.StatusOK, t)
	testNewItemsPage(newItemPosts[:3], t)
}

func TestNewItemAdminPostUsingExistingCompany(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	cookie := loginUser("test1", t)
	client := http.Client{Jar: cookie}
	loadNewItems(client, t)

	cookie = loginUser("admin", t)
	client = http.Client{Jar: cookie}

	currentCompanies, err := db.GetCompanies()
	if err != nil {
		t.Fatal(err)
	}

	pl := postLine{v: url.Values{}}
	pl.newPostLine(4).userID("test1").existingCompany(1).selectToAdd()
	resp, err := client.PostForm(testUrl+"/admin/newitems", pl.v)

	newCompanies, err := db.GetCompanies()
	if len(currentCompanies) != len(newCompanies) {
		t.Fatalf("expected %d items, got %d items", len(currentCompanies), len(newCompanies))
	}

	testPageStatus(resp, err, http.StatusOK, t)
	testNewItemsPage(newItemPosts[:3], t)
}

func TestLimitUserNewItems(t *testing.T) {
}

func TestBadNewItemsPostInvalidAddParams(t *testing.T) {
	server, db := newServer(t)
	defer cleanUp(server, db)
	cookie := loginUser("admin", t)
	client := http.Client{Jar: cookie}
	v := url.Values{"add[]": []string{"abc"}}
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
func (n *newItemHtml) nowt(w string)        {}

func spfi(i int) string { return fmt.Sprintf("%d", i) }

func (p *postLine) newPostLine(itemID int) *postLine {
	p.row = len(p.v["id[]"])
	p.v.Add("id[]", spfi(itemID))
	p.v.Add("isPledge[]", "0")
	p.v.Add("item[]", "0")
	p.v.Add("company[]", "0")
	p.v.Add("userID[]", "")
	p.v.Add("usercompany[]", "")
	p.v.Add("usermake[]", "")
	p.v.Add("usermodel[]", "")
	return p
}

func (p *postLine) selectToAdd() *postLine          { p.v.Add("add[]", spfi(p.row)); return p }
func (p *postLine) selectToDelete() *postLine       { p.v.Add("delete[]", spfi(p.row)); return p }
func (p *postLine) isPledge() *postLine             { p.v["isPledge[]"][p.row] = "1"; return p }
func (p *postLine) userID(u string) *postLine       { p.v["userID[]"][p.row] = u; return p }
func (p *postLine) existingItem(i int) *postLine    { p.v["item[]"][p.row] = spfi(i); return p }
func (p *postLine) existingCompany(c int) *postLine { p.v["company[]"][p.row] = spfi(c); return p }
func (p *postLine) userCompany(c string) *postLine  { p.v["usercompany[]"][p.row] = c; return p }
func (p *postLine) userMake(m string) *postLine     { p.v["usermake[]"][p.row] = m; return p }
func (p *postLine) userModel(m string) *postLine    { p.v["usermodel[]"][p.row] = m; return p }

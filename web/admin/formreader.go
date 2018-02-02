package admin

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
)

// formReader parses a submitted New Items form POST request, captures multiple
// errors that resulted from parsing.
type formReader struct {
	form    url.Values
	row     int
	numRows int
	err     []error
}

func newFormReader(form url.Values) *formReader {
	fr := &formReader{form, -1, len(form["id[]"]), []error{}}
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
	return f.row < f.numRows
}

func (f *formReader) getNewItemPost() newItemPostData {
	if f.errors() {
		panic("getNewItemPost called when formReader in error state")
	}

	i := newItemPostData{
		ID:          f.getInt("id[]"),
		UserID:      f.getString("userID[]"),
		ItemID:      f.getInt("item[]"),
		CompanyID:   f.getInt("company[]"),
		UserCompany: f.getString("usercompany[]"),
		UserMake:    f.getString("usermake[]"),
		UserModel:   f.getString("usermodel[]"),
		Pledge:      f.getString("isPledge[]") == "true",
		Value:       f.getInt("uservalue[]"),
		CurrencyID:  f.getInt("currencyID[]"),
	}
	action := f.getString("action[]")
	log.Printf("action = %+v\n", action)
	switch action {
	case "add":
		i.Add = true
	case "delete":
		i.Delete = true
	case "leave":
		i.Add, i.Delete = false, false
	default:
		f.err = append(f.err, fmt.Errorf("%s is invalid mode in row %d", action, f.row))
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

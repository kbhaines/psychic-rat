package sqldb

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

type mockDB struct {
	nextExpectation int
	expectations    []expectation
	t               *testing.T
}

type expectation struct {
	exec  *execStmt
	query *queryStmt
}

func NewMockDB(t *testing.T) *mockDB {
	return &mockDB{t: t}
}

func (m *mockDB) QueryExpectation(q *queryStmt) *mockDB {
	m.expectations = append(m.expectations, expectation{query: q})
	return m
}

func (m *mockDB) ExecExpectation(e *execStmt) *mockDB {
	m.expectations = append(m.expectations, expectation{exec: e})
	return m
}

func (m *mockDB) CheckAllExpectationsMet() {
	m.t.Helper()
	if m.nextExpectation != len(m.expectations) {
		m.t.Errorf("some interactions missed, expected %d, got %d", len(m.expectations), m.nextExpectation)
	}

}

func (m *mockDB) Exec(query string, args ...interface{}) (Result, error) {
	m.t.Helper()
	if m.nextExpectation == len(m.expectations) {
		m.t.Fatalf("no more expectations set for Exec")
	}
	exec := m.expectations[m.nextExpectation]
	m.nextExpectation++
	if exec.exec != nil {
		checkExecInsert(m.t, exec.exec, query, args)
		return exec.exec, nil
	}
	return exec.exec, fmt.Errorf("not able to match mock exec for query %s", query)
}

func (m *mockDB) Query(query string, args ...interface{}) (Rows, error) {
	m.t.Helper()
	if m.nextExpectation == len(m.expectations) {
		m.t.Fatalf("no more expectations set for Query")
	}
	exec := m.expectations[m.nextExpectation]
	m.nextExpectation++
	if exec.query != nil {
		checkQuery(m.t, exec.query, query, args)
		if exec.query.err != nil {
			return exec.query.rows, exec.query.err
		}
		return exec.query.rows, nil
	}
	return nil, fmt.Errorf("not able to match mock to query %s", query)
}

func (m *mockDB) QueryRow(query string, args ...interface{}) Row {
	m.t.Helper()
	if m.nextExpectation == len(m.expectations) {
		m.t.Fatalf("no more expectations set for QueryRow")
	}
	exec := m.expectations[m.nextExpectation]
	m.nextExpectation++
	if exec.query != nil {
		if exec.query.err != nil {
			return &rowError{err: exec.query.err}
		}
		if len(exec.query.rows.rows) != 1 {
			m.t.Errorf("QueryRow called, expected 1 row in mock results, got %d", len(exec.query.rows.rows))
		}
		checkQuery(m.t, exec.query, query, args)
		return exec.query.rows
	}
	return &rowError{query: query}
}

func (m mockDB) Close() {}

// execStmt represents an expectation for an sql.Exec call

type execStmt struct {
	table        string
	columns      map[string]interface{}
	insertId     int64
	rowsAffected int64
}

func NewExec(table string) *execStmt {
	return &execStmt{table: table, columns: make(map[string]interface{})}
}

func (e *execStmt) WithColumnValue(col string, value interface{}) *execStmt {
	e.columns[col] = value
	return e
}

func (e *execStmt) WithInsertId(id int64) *execStmt {
	e.insertId = id
	return e
}

func (e *execStmt) LastInsertId() (int64, error) {
	return e.insertId, nil
}

// queryStmt represents an expectation for a sql.Query call
type queryStmt struct {
	table       string
	columns     string
	whereClause string
	rows        *rowsResult
	err         error
}

func NewQuery(table string) *queryStmt {
	return &queryStmt{table: table, rows: &rowsResult{}, err: nil}
}

func (q *queryStmt) WithColumns(cols ...string) *queryStmt {
	q.columns = strings.Join(cols, ",")
	return q
}

func (q *queryStmt) WithError(msg string) *queryStmt {
	q.err = fmt.Errorf(msg)
	return q
}

func (q *queryStmt) WithResultsRow(v ...interface{}) *queryStmt {
	q.rows.rows = append(q.rows.rows, v)
	return q
}

// rowsResult represents a result from a sql.Query call

type rowsResult struct {
	next int
	rows [][]interface{}
}

func (r *rowsResult) Next() bool {
	return r.next < len(r.rows)
}

func (r *rowsResult) Scan(v ...interface{}) error {
	if r.next == len(r.rows) {
		// TODO: make test fail
		return fmt.Errorf("Scan called but rows exhausted")
	}
	row := r.rows[r.next]
	if len(row) != len(v) {
		// TODO: make test fail
		return fmt.Errorf("Scan called with wrong number of arguments, expected %d got %d", len(row), len(v))
	}

	for i, val := range row {
		// TODO: make test fail
		if err := convertAssign(v[i], val); err != nil {
			return err
		}
	}
	r.next++

	return nil
}

// rowError is the case where a sql.QueryRow doesn't return just a single row from
// the mock.
type rowError struct {
	query string
	err   error
}

func (r *rowError) Scan(v ...interface{}) error {
	if r.err != nil {
		return r.err
	}
	return fmt.Errorf("not able to match mock to query %s", r.query)
}

/////

func checkExecInsert(t *testing.T, insert *execStmt, query string, args []interface{}) {
	t.Helper()

	re := regexp.MustCompile("insert\\s+into\\s+(.+)\\s*\\((.*)\\)\\s+values\\s*\\((.*)\\)")
	results := re.FindStringSubmatch(query)
	if len(results) != 4 {
		t.Fatalf("could not match exec statement; regexp match [%v] from [%v]", results, query)
	}

	table := results[1]
	columns := strings.Split(strings.Replace(results[2], " ", "", -1), ",")
	values := strings.Split(strings.Replace(results[3], " ", "", -1), ",")

	if insert.table != table {
		t.Fatalf("wrong table, expected %v, got %v", insert.table, table)
	}
	if len(columns) != len(values) {
		t.Fatalf("wrong number of values, expected %v, got %v", len(columns), len(values))
	}
	if len(insert.columns) != len(columns) {
		t.Fatalf("column counts don't match, expected %v, got %v", insert.columns, columns)
	}
	for i, col := range columns {
		expv, exists := insert.columns[col]
		if !exists {
			t.Fatalf("unexpected column: %v", col)
		}
		if values[i] != "?" {
			t.Fatalf("got unexpected placeholder, expected ?, got %v", values[i])
		}
		if !reflect.DeepEqual(expv, args[i]) {
			t.Fatalf("types & values don't match, expected %v of type %T, got %v of type %T", expv, expv, args[i], args[i])
		}
	}
}

func checkQuery(t *testing.T, expQuery *queryStmt, query string, args []interface{}) {
	t.Helper()
	re := regexp.MustCompile("select\\s+(.*)\\s+from\\s+(\\w*)(\\s+where\\s+(.*))?")
	match := re.FindStringSubmatch(query)
	if len(match) < 3 {
		t.Fatalf("could not match select statement (%v)", query)
	}
	columns := strings.Replace(match[1], " ", "", -1)
	table := match[2]
	if expQuery.table != table {
		t.Fatalf("expected table %s, got %s in query %s", expQuery.table, table, query)
	}
	if expQuery.columns != columns {
		t.Fatalf("expected columns [%s], got [%s] in query %s", expQuery.columns, columns, query)
	}
}

////////////////////////////////////////////////////////////////////////////////
// The following lines are taken from database/sql/convert.go

var errNilPtr = fmt.Errorf("nil pointer")

type RawBytes []byte

// convertAssign copies to dest the value in src, converting it if possible.
// An error is returned if the copy would result in loss of information.
// dest should be a pointer type.
func convertAssign(dest, src interface{}) error {
	// Common cases, without reflect.
	switch s := src.(type) {
	case string:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return errNilPtr
			}
			*d = s
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = []byte(s)
			return nil
		case *RawBytes:
			if d == nil {
				return errNilPtr
			}
			*d = append((*d)[:0], s...)
			return nil
		}
	case []byte:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return errNilPtr
			}
			*d = string(s)
			return nil
		case *interface{}:
			if d == nil {
				return errNilPtr
			}
			*d = cloneBytes(s)
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = cloneBytes(s)
			return nil
		case *RawBytes:
			if d == nil {
				return errNilPtr
			}
			*d = s
			return nil
		}
	case time.Time:
		switch d := dest.(type) {
		case *time.Time:
			*d = s
			return nil
		case *string:
			*d = s.Format(time.RFC3339Nano)
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = []byte(s.Format(time.RFC3339Nano))
			return nil
		case *RawBytes:
			if d == nil {
				return errNilPtr
			}
			*d = s.AppendFormat((*d)[:0], time.RFC3339Nano)
			return nil
		}
	case nil:
		switch d := dest.(type) {
		case *interface{}:
			if d == nil {
				return errNilPtr
			}
			*d = nil
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = nil
			return nil
		case *RawBytes:
			if d == nil {
				return errNilPtr
			}
			*d = nil
			return nil
		}
	}

	var sv reflect.Value

	switch d := dest.(type) {
	case *string:
		sv = reflect.ValueOf(src)
		switch sv.Kind() {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			*d = asString(src)
			return nil
		}
	case *[]byte:
		sv = reflect.ValueOf(src)
		if b, ok := asBytes(nil, sv); ok {
			*d = b
			return nil
		}
	case *RawBytes:
		sv = reflect.ValueOf(src)
		if b, ok := asBytes([]byte(*d)[:0], sv); ok {
			*d = RawBytes(b)
			return nil
		}
	case *bool:
		bv, err := driver.Bool.ConvertValue(src)
		if err == nil {
			*d = bv.(bool)
		}
		return err
	case *interface{}:
		*d = src
		return nil
	}

	dpv := reflect.ValueOf(dest)
	if dpv.Kind() != reflect.Ptr {
		return errors.New("destination not a pointer")
	}
	if dpv.IsNil() {
		return errNilPtr
	}

	if !sv.IsValid() {
		sv = reflect.ValueOf(src)
	}

	dv := reflect.Indirect(dpv)
	if sv.IsValid() && sv.Type().AssignableTo(dv.Type()) {
		switch b := src.(type) {
		case []byte:
			dv.Set(reflect.ValueOf(cloneBytes(b)))
		default:
			dv.Set(sv)
		}
		return nil
	}

	if dv.Kind() == sv.Kind() && sv.Type().ConvertibleTo(dv.Type()) {
		dv.Set(sv.Convert(dv.Type()))
		return nil
	}

	// The following conversions use a string value as an intermediate representation
	// to convert between various numeric types.
	//
	// This also allows scanning into user defined types such as "type Int int64".
	// For symmetry, also check for string destination types.
	switch dv.Kind() {
	case reflect.Ptr:
		if src == nil {
			dv.Set(reflect.Zero(dv.Type()))
			return nil
		}
		dv.Set(reflect.New(dv.Type().Elem()))
		return convertAssign(dv.Interface(), src)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s := asString(src)
		i64, err := strconv.ParseInt(s, 10, dv.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetInt(i64)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s := asString(src)
		u64, err := strconv.ParseUint(s, 10, dv.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetUint(u64)
		return nil
	case reflect.Float32, reflect.Float64:
		s := asString(src)
		f64, err := strconv.ParseFloat(s, dv.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetFloat(f64)
		return nil
	case reflect.String:
		switch v := src.(type) {
		case string:
			dv.SetString(v)
			return nil
		case []byte:
			dv.SetString(string(v))
			return nil
		}
	}

	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", src, dest)
}

func cloneBytes(b []byte) []byte {
	if b == nil {
		return nil
	}
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func asString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}
	return fmt.Sprintf("%v", src)
}

func asBytes(buf []byte, rv reflect.Value) (b []byte, ok bool) {
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.AppendInt(buf, rv.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.AppendUint(buf, rv.Uint(), 10), true
	case reflect.Float32:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 32), true
	case reflect.Float64:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 64), true
	case reflect.Bool:
		return strconv.AppendBool(buf, rv.Bool()), true
	case reflect.String:
		s := rv.String()
		return append(buf, s...), true
	}
	return
}

func strconvErr(err error) error {
	if ne, ok := err.(*strconv.NumError); ok {
		return ne.Err
	}
	return err
}

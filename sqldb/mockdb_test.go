package sqldb

import (
	"fmt"
	"reflect"
	"regexp"
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
	exec := m.expectations[m.nextExpectation]
	m.nextExpectation++
	if exec.exec != nil {
		checkExecInsert(m.t, exec.exec, query, args)
		return exec.exec, nil
	}
	return exec.exec, fmt.Errorf("not able to match mock exec for query %s", query)
}

func (m *mockDB) Query(query string, args ...interface{}) (Rows, error) {
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

	re := regexp.MustCompile("insert into (.*)\\((.*)\\) values\\((.*)\\)")
	results := re.FindStringSubmatch(query)
	if len(results) != 4 {
		t.Fatalf("could not match query: %v", query)
	}

	table := results[1]
	columns := strings.Split(results[2], ",")
	values := strings.Split(results[3], ",")

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
	re := regexp.MustCompile("select (.*) from (\\w*)( where (.*))?")
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

var errNilPtr = fmt.Errorf("nil pointer")

// convertAssign copies to dest the value in src, converting it if possible.
// An error is returned if the copy would result in loss of information.
// dest should be a pointer type.
func convertAssign(dest, src interface{}) error {
	// Common cases, without reflect.
	switch s := src.(type) {
	case int:
		switch d := dest.(type) {
		case *int:
			if d == nil {
				return errNilPtr
			}
			*d = s
			return nil
		}
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
		}
	case []byte:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return errNilPtr
			}
			*d = string(s)
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
		}
	case float64:
		switch d := dest.(type) {
		case *float64:
			*d = s
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
		}
	}

	return fmt.Errorf("unable to parse in Scan")
}

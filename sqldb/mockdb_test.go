package sqldb

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

type expectedExecStmt struct {
	insert       *insertStmt
	query        *queryStmt
	insertId     int64
	rowsAffected int64
}

type insertStmt struct {
	table   string
	columns map[string]interface{}
}

type queryStmt struct {
	table       string
	columns     string
	whereClause string
	rows        *rowsResult
}

type rowsResult struct {
	next int
	rows [][]interface{}
}

type mockDB struct {
	execsDone int
	execs     []expectedExecStmt
	t         *testing.T
}

func NewQuery(table string) *queryStmt {
	return &queryStmt{table: table, rows: &rowsResult{}}
}

func (q *queryStmt) WithColumns(cols ...string) *queryStmt {
	q.columns = strings.Join(cols, ",")
	return q
}

func (q *queryStmt) WithResultsRow(v ...interface{}) *queryStmt {
	q.rows.rows = append(q.rows.rows, v)
	return q
}

func (m expectedExecStmt) LastInsertId() (int64, error) {
	return m.insertId, nil
}

func (m expectedExecStmt) RowsAffected() (int64, error) {
	return m.rowsAffected, nil
}

func (m mockDB) Exec(query string, args ...interface{}) (Result, error) {
	exec := m.execs[m.execsDone]
	m.execsDone++
	if exec.insert != nil {
		checkExecInsert(m.t, exec.insert, query, args)
		return exec, nil
	}
	return exec, fmt.Errorf("not able to match mock exec for query %s", query)
}

func (m mockDB) Query(query string, args ...interface{}) (Rows, error) {
	exec := m.execs[m.execsDone]
	m.execsDone++
	if exec.query != nil {
		checkQuery(m.t, exec.query, query, args)
		return exec.query.rows, nil
	}
	return nil, fmt.Errorf("not able to match mock to query %s", query)
}

func (m mockDB) Close() {
}

func checkExecInsert(t *testing.T, insert *insertStmt, query string, args []interface{}) {
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
	re := regexp.MustCompile("select (.*) from (.*)( where (.*))?")
	match := re.FindStringSubmatch(query)
	if len(match) < 3 {
		t.Fatalf("could not match select statement (%v)", query)
	}
	columns := match[1]
	table := match[2]
	if expQuery.table != table {
		t.Fatalf("expected table %s, got %s in query %s", expQuery.table, table, query)
	}
	if expQuery.columns != columns {
		t.Fatalf("expected columns %s, got %s in query %s", expQuery.columns, columns, query)
	}
}

func (m mockDB) QueryRow(query string, args ...interface{}) Row {
	panic("not implemented")
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

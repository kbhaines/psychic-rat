package sqldb

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
)

type SQLLite3 struct {
	DBInterface
	sqllite *sql.DB
}

func NewSqlLiteDB(name string) (*SQLLite3, error) {
	if _, err := os.Stat(name); os.IsExist(err) {
		panic(fmt.Sprintf("refusing to create when DB %s already exists", name))
	}
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}
	return setupDB(db)
}

func OpenDB(name string) (*SQLLite3, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}
	return setupDB(db)
}

func Backup(originalFile, backupFile string) error {
	os.Remove(backupFile)
	return exec.Command("sqlite3", originalFile, ".backup "+backupFile).Run()
}

func setupDB(db *sql.DB) (*SQLLite3, error) {
	schemaUpdates(db)
	_, err := db.Exec("PRAGMA synchronous = OFF")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &SQLLite3{sqllite: db}, nil
}

func (s *SQLLite3) Exec(query string, args ...interface{}) (Result, error) {
	panic("not implemented")
}

func (s *SQLLite3) Query(query string, args ...interface{}) (Rows, error) {
	panic("not implemented")
}

func (s *SQLLite3) QueryRow(query string, args ...interface{}) Row {
	panic("not implemented")
}

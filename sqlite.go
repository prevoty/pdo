// +build cgo

package pdo

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type (
	Sqlite struct {
		DBO
	}
)

func NewSqlite(dsn string) (*Sqlite, error) {

	_, err := os.Stat(dsn)
	if err != nil {
		return nil, err
	}

	s := new(Sqlite)

	s.DB, err = sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	err = s.DB.Ping()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s (%s)\n", err, dsn))
	}

	return s, nil

}

func (s *Sqlite) StartTransaction() error {
	_, err := s.DB.Exec("BEGIN TRANSACTION")
	if err != nil {
		return err
	}
	return nil
}
func (s *Sqlite) Rollback() error {
	_, err := s.DB.Exec("ROLLBACK")
	if err != nil {
		return err
	}
	return nil
}
func (s *Sqlite) Commit() error {
	_, err := s.DB.Exec("COMMIT")
	if err != nil {
		return err
	}
	return nil
}

// +build cgo

package pdo

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
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

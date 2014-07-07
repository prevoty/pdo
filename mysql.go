package pdo

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type (
	MySQL struct {
		DBO
	}
)

func NewMySQL(dsn string) (*MySQL, error) {

	var (
		m   = new(MySQL)
		err error
	)

	m.DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = m.DB.Ping()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s (%s)\n", err, dsn))
	}

	return m, nil

}

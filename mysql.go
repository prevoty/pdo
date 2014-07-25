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

func (m *MySQL) StartTransaction() error {
	_, err := m.DB.Exec("START TRANSACTION")
	if err != nil {
		return err
	}
	return nil
}
func (m *MySQL) Rollback() error {
	_, err := m.DB.Exec("ROLLBACK")
	if err != nil {
		return err
	}
	return nil
}
func (m *MySQL) Commit() error {
	_, err := m.DB.Exec("COMMIT")
	if err != nil {
		return err
	}
	return nil
}

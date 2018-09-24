package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type ConnWrapper struct {
	*sqlx.DB
}

//go:generate counterfeiter -o fakes/transaction.go --fake-name Transaction . Transaction
type Transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	Commit() error
	Rollback() error
	Rebind(string) string
	DriverName() string
}

func (c *ConnWrapper) Beginx() (Transaction, error) {
	return c.DB.Beginx()
}

func (c *ConnWrapper) OpenConnections() int {
	return c.DB.Stats().OpenConnections
}

func (c *ConnWrapper) RawConnection() *sqlx.DB {
	return c.DB
}

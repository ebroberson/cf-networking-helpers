package db

import (
	"code.cloudfoundry.org/bbs/db/sqldb/helpers/monitor"
	"github.com/jmoiron/sqlx"
)

type ConnWrapper struct {
	*sqlx.DB
	Monitor monitor.Monitor
}

func (c *ConnWrapper) Beginx() (Transaction, error) {
	var innerTx *sqlx.Tx
	err := c.Monitor.Monitor(func() error {
		var err error
		innerTx, err = c.DB.Beginx()
		return err
	})

	tx := &monitoredTx{
		tx:      innerTx,
		monitor: c.Monitor,
	}

	return tx, err
}

func (c *ConnWrapper) OpenConnections() int {
	return c.DB.Stats().OpenConnections
}

func (c *ConnWrapper) RawConnection() *sqlx.DB {
	return c.DB
}

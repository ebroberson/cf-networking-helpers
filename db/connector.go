package db

import (
	"database/sql"
	"fmt"
	"net"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type RetriableError struct {
	Inner error
	Msg   string
}

func (r RetriableError) Error() string {
	return fmt.Sprintf("%s: %s", r.Msg, r.Inner.Error())
}

func GetConnectionPool(dbConfig Config) (*sqlx.DB, error) {
	connectionString, err := dbConfig.ConnectionString()
	if err != nil {
		panic(err)
	}
	nativeDBConn, err := sql.Open(dbConfig.Type, connectionString)
	if err != nil {
		return nil, fmt.Errorf("unable to open database connection: %s", err)
	}

	dbConn := sqlx.NewDb(nativeDBConn, dbConfig.Type)

	if err = dbConn.Ping(); err != nil {
		dbConn.Close()
		if netErr, ok := err.(*net.OpError); ok {
			return nil, RetriableError{
				Inner: netErr,
				Msg:   "unable to ping",
			}
		}
		return nil, fmt.Errorf("unable to ping: %s", err)
	}

	return dbConn, nil
}

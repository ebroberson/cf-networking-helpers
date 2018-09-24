package db

import (
	"time"

	"fmt"
)

//go:generate counterfeiter -o ../fakes/sleeper.go --fake-name Sleeper . sleeper
type sleeper interface {
	Sleep(time.Duration)
}

type SleeperFunc func(time.Duration)

func (sf SleeperFunc) Sleep(duration time.Duration) {
	sf(duration)
}

type RetriableConnector struct {
	Connector     func(Config) (*ConnWrapper, error)
	Sleeper       sleeper
	RetryInterval time.Duration
	MaxRetries    int
}

func (r *RetriableConnector) GetConnectionPool(dbConfig Config) (*ConnWrapper, error) {
	var attempts int
	for {
		attempts++

		db, err := r.Connector(dbConfig)
		if err == nil {
			return db, nil
		}

		if _, ok := err.(RetriableError); ok && attempts < r.MaxRetries {
			println(fmt.Sprintf("retrying due to getting an error %#+v", err))
			r.Sleeper.Sleep(r.RetryInterval)
			continue
		}

		return nil, err
	}
}

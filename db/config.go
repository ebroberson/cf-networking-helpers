package db

import (
	"fmt"
	"time"
	"io/ioutil"
	"crypto/x509"
	"crypto/tls"
	"github.com/go-sql-driver/mysql"
)

type Config struct {
	Type         string `json:"type" validate:"nonzero"`
	User         string `json:"user" validate:"nonzero"`
	Password     string `json:"password"`
	Host         string `json:"host" validate:"nonzero"`
	Port         uint16 `json:"port" validate:"nonzero"`
	Timeout      int    `json:"timeout" validate:"min=1"`
	DatabaseName string `json:"database_name" validate:""`
	RequireSSL   bool   `json:"require_ssl" validate:""`
	CACert       string `json:"ca_cert" validate:""`
}

func (c Config) ConnectionString() (string, error) {
	if c.Timeout < 1 {
		return "", fmt.Errorf("timeout must be at least 1 second: %d", c.Timeout)
	}
	switch c.Type {
	case "postgres":
		ms := (time.Duration(c.Timeout) * time.Second).Nanoseconds() / 1000 / 1000
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&connect_timeout=%d", c.User, c.Password, c.Host, c.Port, c.DatabaseName, ms), nil
	case "mysql":
		return c.buildMysqlConnectionString()
	default:
		return "", fmt.Errorf("database type '%s' is not supported", c.Type)
	}
}

func (c Config) buildMysqlConnectionString() (string, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", c.User, c.Password, c.Host, c.Port, c.DatabaseName)
	dbConfig, err := mysql.ParseDSN(connString)
	if err != nil {
		return "", fmt.Errorf("parsing db connection string: %s", err)
	}

	timeoutDuration := time.Duration(c.Timeout) * time.Second
	dbConfig.Timeout = timeoutDuration
	dbConfig.ReadTimeout = timeoutDuration
	dbConfig.WriteTimeout = timeoutDuration

	if c.RequireSSL {
		certBytes, err := ioutil.ReadFile(c.CACert)
		if err != nil {
			return "", fmt.Errorf("reading db ca cert file: %s", err)
		}

		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(certBytes); !ok {
			return "", fmt.Errorf("appending cert to pool from pem - invalid cert bytes")
		}

		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			RootCAs:            caCertPool,
		}

		tlsConfigName := fmt.Sprintf("%s-tls", c.DatabaseName)

		err = mysql.RegisterTLSConfig(tlsConfigName, tlsConfig)
		if err != nil {
			return "", fmt.Errorf("registering mysql tls config: %s", err)
		}

		dbConfig.TLSConfig = tlsConfigName
	}

	return dbConfig.FormatDSN(), nil
}

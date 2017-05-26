package db_test

import (
	"code.cloudfoundry.org/cf-networking-helpers/db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var (
		config db.Config
	)
	BeforeEach(func() {
		config = db.Config{
			User:         "some-user",
			Password:     "some-password",
			Host:         "some-host",
			Port:         uint16(1234),
			DatabaseName: "some-database",
			Timeout:      5,
		}
	})
	Describe("ConnectionString", func() {
		Context("when the type is postgres", func() {
			BeforeEach(func() {
				config.Type = "postgres"
			})
			It("returns the connection string", func() {
				connectionString, err := config.ConnectionString()
				Expect(err).NotTo(HaveOccurred())
				Expect(connectionString).To(Equal("postgres://some-user:some-password@some-host:1234/some-database?sslmode=disable&connect_timeout=5000&read_timeout=5000&write_timeout=5000"))
			})
		})
		Context("when the type is mysql", func() {
			BeforeEach(func() {
				config.Type = "mysql"
			})
			It("returns the connection string", func() {
				connectionString, err := config.ConnectionString()
				Expect(err).NotTo(HaveOccurred())
				Expect(connectionString).To(Equal("some-user:some-password@tcp(some-host:1234)/some-database?parseTime=true&timeout=5s&readTimeout=5s&writeTimeout=5s"))
			})
		})

		Context("when the type is neither", func() {
			BeforeEach(func() {
				config.Type = "neither"
			})
			It("returns an error", func() {
				_, err := config.ConnectionString()
				Expect(err).To(MatchError("database type 'neither' is not supported"))
			})
		})

		Context("when the timeout is less than 1", func() {
			BeforeEach(func() {
				config.Type = "postgres"
				config.Timeout = 0
			})
			It("returns an error", func() {
				_, err := config.ConnectionString()
				Expect(err).To(MatchError("timeout must be at least 1 second: 0"))
			})
		})
	})

})

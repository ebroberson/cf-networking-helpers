package db_test

import (
	"io/ioutil"

	"code.cloudfoundry.org/cf-networking-helpers/db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/go-sql-driver/mysql"
)

const (
	DATABASE_CA_CERT = `-----BEGIN CERTIFICATE-----
MIIE4jCCAsqgAwIBAgIBATANBgkqhkiG9w0BAQsFADARMQ8wDQYDVQQDEwZmYWtl
Q0EwHhcNMTgwNTEwMjM1MDM2WhcNMTkxMTEwMjM1MDM2WjARMQ8wDQYDVQQDEwZm
YWtlQ0EwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQC3u/I7qztSp8rO
S266wo53NtqdtM/8iyyCigqCwHgJ7CauGKq33zTBaUkRljeRn/AXvkChPPEA3KQm
Wrv5YkFhCq/1EOB2JUMPVrUXjP/g6HwPAHX2IvC0pJoYMmb7TloGXfTjV/A/2e41
Q1zRSWAMDXCUfnAn6skkihV9YGipdM/r0+O9n8tb3F+Z+dYvMu89DwyptI/oNzNK
DyzkQf5WZ1PCqEow7ZcbSQP3RH2Ds6I+AG98nxB4irsmUkoZnUhQzTc9DpINmgI8
3Yg8YyTFODZ1BbsnST1Y01rWMvVkXy/89+fBqN4kGI12CYtbh69Shr/0cim8fETT
N9CLzqDpPlnfAJGv/VwVSzwxhYuYBfh3PtlAo5OfVBhYoGq9npjV3H/j9N9r0aE4
MkQvYaATB/fQ823mtjLDqtkIvZCXq1PZA90oQ87n1FPeklc9/T14SXNcHVMBMpSX
mPdaJvBoXjlwl1EKvZIQzz/luxMZfgqSRy4TLcJKJ+E+3bU2RZcz56r5aTV8+9aS
/SL80oQpGzXK4pWcFvELlGcW2LnP7XPE3t1HzS0kEVSFnVyw4/UJvcsyZSUl2bDs
FJl0HOkVuNtjnhCKTiRpRTYdKxKhvxp46/0FMtnujIq1WQF6yKUK3mOUpJpaJe+4
3fd7UsX8Qz0Gwj9scxBCTTNeVXkU8QIDAQABo0UwQzAOBgNVHQ8BAf8EBAMCAQYw
EgYDVR0TAQH/BAgwBgEB/wIBADAdBgNVHQ4EFgQUsUmpQ5m8K0hqWlw/7M5Q1wG0
FO4wDQYJKoZIhvcNAQELBQADggIBAKWrPjCEYWWWnWWjIFazsbe98eSu7N4yTDtc
yN4D+k0sIWbbgKbeAV9k5N5H8p9w5tNzjsjUCK88qdEvN+0kJHWCvt2zffBlP6tX
nC2bB12CjPPYIUpnG14ghZB/Uxj240eo/1JCrsb/qTecW2H3UbLmjtmx10RJVP7U
kseGnsXQwPIEgOVHubVLkIobv88zLSJKgf8syhnbihl5/eKIBMfreaI7mW2+CqTE
Y1SfIPTpU59YHW7TNLsI9WgQNtDqCORKwzzpVnUWPfO0iQ1+wEnjfhsMDOmzfBLO
l1HspfZpRnOXZFROnuNvR1V+qyPrKMm/F01B7Z4ESxa7ktNEbrwOt2wGT/keXocw
z1LXbrG/WBjov4HCD6pXv+w4XwkR9bPEHkrMZ/INCm5oIq7JLTZcjb56aawqAk8W
0XXKhjFTIGO46GPTbcJTWxs3BJX8C2mL5aVHWekJfuXsCU/0GIxudV8VQJqulq/1
dlZjOpycEZ11hWkZENsJ8ddDX0eYTR95MGAq8J7m1Q0Ts0X/d5ATc1mREf3wqhSn
TFFl82cBZE15vJfk5ekNof8Hx2NTZYwfplKKqb8epo2pIA/j3/PRjo80AFyicjoY
7/Xiu2K2JGmsEF3XQVowXVsxngkLZSqHml+WRweqaK48Zbojj/hUkz+xOAoucqZO
VCbEyl6T
-----END CERTIFICATE-----`
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
				Expect(connectionString).To(Equal("postgres://some-user:some-password@some-host:1234/some-database?sslmode=disable&connect_timeout=5000"))
			})
		})

		Context("when the type is mysql", func() {
			BeforeEach(func() {
				config.Type = "mysql"
			})

			It("returns the connection string", func() {
				connectionString, err := config.ConnectionString()
				Expect(err).NotTo(HaveOccurred())
				Expect(connectionString).To(Equal("some-user:some-password@tcp(some-host:1234)/some-database?parseTime=true&readTimeout=5s&timeout=5s&writeTimeout=5s"))
			})

			Context("when require_ssl is enabled", func() {
				BeforeEach(func() {
					config.RequireSSL = true
				})

				AfterEach(func() {
					mysql.DeregisterTLSConfig("some-database-tls")
				})

				Context("success", func() {
					BeforeEach(func() {
						caCertFile, err := ioutil.TempFile("", "")
						Expect(err).NotTo(HaveOccurred())

						_, err = caCertFile.Write([]byte(DATABASE_CA_CERT))
						Expect(err).NotTo(HaveOccurred())

						config.CACert = caCertFile.Name()
					})

					It("returns the amended connection string", func() {
						connectionString, err := config.ConnectionString()
						Expect(err).NotTo(HaveOccurred())
						Expect(connectionString).To(Equal("some-user:some-password@tcp(some-host:1234)/some-database?parseTime=true&readTimeout=5s&timeout=5s&tls=some-database-tls&writeTimeout=5s"))
					})
				})

				Context("when reading the cert file fails", func() {
					BeforeEach(func() {
						config.CACert = "garbage"
					})

					It("returns an error", func() {
						_, err := config.ConnectionString()
						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError("reading db ca cert file: open garbage: no such file or directory"))
					})
				})

				Context("when adding the cert to the pool fails", func() {
					BeforeEach(func() {
						caCertFile, err := ioutil.TempFile("", "")
						Expect(err).NotTo(HaveOccurred())

						_, err = caCertFile.Write([]byte("garbage"))
						Expect(err).NotTo(HaveOccurred())

						config.CACert = caCertFile.Name()
					})


					It("returns an error", func() {
						_, err := config.ConnectionString()
						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError("appending cert to pool from pem - invalid cert bytes"))
					})
				})
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

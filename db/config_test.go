package db_test

import (
	"io/ioutil"

	"code.cloudfoundry.org/cf-networking-helpers/db"

	"net/url"

	"github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	DATABASE_CA_CERT = `-----BEGIN CERTIFICATE-----
MIIE6DCCAtCgAwIBAgIBATANBgkqhkiG9w0BAQsFADATMREwDwYDVQQDEwhDZXJ0
QXV0aDAgFw0yMTAzMDEyMzQ4MjFaGA8yMTIxMDMwMTIzNDgxNlowEzERMA8GA1UE
AxMIQ2VydEF1dGgwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDObRYi
LqJ0B4yA71FPv+ydmTKiOzTPsL7GXi2PwdA+J6y7t7yW3PK1U3KPLcOd1uMxe/mj
jVS7NhZ8b2S6ROhY2yhbGhH5gZOh4VQKOb7K54iyNdryeoIPgI2lxn6Nxa+MSxPM
b1iHZMCauWkr0GqyP59zoEs+b1KNvnPPi39QySyZH+eTptZUcc1PFvwpbA31oJP7
NIIl7jOpvg6wRz04OtCyYuyoyWlJ9N4BnMMxdIhnNGozUm0pJglkDGFHAMREG4Ah
ZaVCuib3w8nPA/cv2BAYY1uQwCsUjA9/WUCct7+wVw9pzKSJXc5DEr8r+8sAb4Zi
X0GmElAS/Qb0OzkNSboy1iI/+QyfLPEhOsWXinqFFtuhuXP6Ek4YslMzlRGQMYjb
0BDypF9DDcLHKQyE7EKuxILDstQyucniea+Z+c0+CUtNhdiuG61o5yjBRO4ELffX
xoSyk2hoinPZBh3zrCF4Em4wScPkgYdcQjwRVuBdBljZawTDF8kSIuaCVoCU0NdC
MX60TTYzG4ij/2XQ5bDuvyPMyQJBu5PzXU+jPDYO/ZuA0EW+VVdjm8wnbmYI+b/f
aLQL6gHTWHvYDBVWgdiyPi59sMgCMM6JgvIXKy66Ghj7lyMHxNRHqHCb2T4WarK/
wPIuBMfpcfOZnVuWFxBRQNIj1oqkKM3LiYnDBwIDAQABo0UwQzAOBgNVHQ8BAf8E
BAMCAQYwEgYDVR0TAQH/BAgwBgEB/wIBADAdBgNVHQ4EFgQUH/3zqYnc5/64PR+U
fIdox3SOntQwDQYJKoZIhvcNAQELBQADggIBAMpN1O+DwbwkxqfG9wlPWEZ0fmN4
FveMUVc28p5Bh+F6jEKj8emut2Np5mQnFShN3FzvXoCdKCgcK9rV/YMkGfpO2MKk
vrKzd6m5GG1WSUVuAvYWsJP4TAh+nMwb7tF3NRyVOQE5RxfI2ahDbyRzUpj77f9+
UPZ3rO51BxfrIgpeKan6Hngqcz3p+ytXiyb19I/oz5HkkjWkuZLdFWlbrJ3M24vT
G32fz/rXeS44jlJwtfnnX+a9WfsTv7DUbcWmkWXtnl+k/w2RLmmeu61URktR3Tz6
a7jaDarP8jxUhJYfrmGNpiRfHUR7Rgg7G25832D3iH5foSPugaabFgAgHyOG8zN5
8P8JY6JQlU7vURbZSFgD3FPE8Df9tAevRVW2HlZ5kMFsiF7vxX9mexCP4Akw4E4N
G8PlCj1p+AFk3ZEDDEQlZWOICNgOm7SozIpkCTt9p9kDJ468ukUZxLi/TVygp1BC
j6A2/SdSFxwuYj/LBw7t7czYrBqL9yDKuRLHmDvXrSpSI7XtHqK0ZUYbXiWlKsVg
xgyCdeHzYNblscbd/915lUcrhV671hSY5W70WpERUD2HBKvCwVQcTdLizX2bC/sZ
+zVVJD8zh6r+HCKfmN96wcZPsgWDHWNCXEMg2T8Tmj/MG7majgz94WfqefpjoVi2
x8WmHPsvkb8pdlch
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
				Expect(connectionString).To(Equal("postgres://some-user:some-password@some-host:1234/some-database?connect_timeout=5000&sslmode=disable"))
			})

			Context("when ssl is required", func() {
				BeforeEach(func() {
					config.RequireSSL = true
					config.CACert = "/tmp/cert"
				})

				Context("when skip_hostname_validation is set", func() {
					BeforeEach(func() {
						config.SkipHostnameValidation = true
					})
					It("sets sslmode to \"require\"", func() {
						connectionString, err := config.ConnectionString()
						Expect(err).NotTo(HaveOccurred())
						connUrl, err := url.Parse(connectionString)
						Expect(err).NotTo(HaveOccurred())
						connQuery := connUrl.Query()
						Expect(connQuery.Get("sslmode")).To(Equal("require"))
					})
				})

				Context("when skip_hostname_validation is not set", func() {
					Context("when ca_cert is empty", func() {
						BeforeEach(func() {
							config.CACert = ""
						})
						It("returns an error", func() {
							_, err := config.ConnectionString()
							Expect(err).To(HaveOccurred())
						})
					})

					It("sets sslmode to \"verify-full\"", func() {
						connectionString, err := config.ConnectionString()
						Expect(err).NotTo(HaveOccurred())
						connUrl, err := url.Parse(connectionString)
						Expect(err).NotTo(HaveOccurred())
						connQuery := connUrl.Query()
						Expect(connQuery.Get("sslmode")).To(Equal("verify-full"))
					})
					It("sets sslrootcert", func() {
						connectionString, err := config.ConnectionString()
						Expect(err).NotTo(HaveOccurred())
						connUrl, err := url.Parse(connectionString)
						Expect(err).NotTo(HaveOccurred())
						connQuery := connUrl.Query()
						Expect(connQuery.Get("sslrootcert")).To(Equal("/tmp/cert"))
					})
				})

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

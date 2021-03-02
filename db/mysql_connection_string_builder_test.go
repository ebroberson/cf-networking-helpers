package db_test

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"

	"code.cloudfoundry.org/cf-networking-helpers/db"
	"code.cloudfoundry.org/cf-networking-helpers/fakes"
	"github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	DATABASE_CLIENT_CERT = `-----BEGIN CERTIFICATE-----
MIIEIDCCAgigAwIBAgIQFs8G9Je7Qf7pQCvt0Mq3qTANBgkqhkiG9w0BAQsFADAT
MREwDwYDVQQDEwhDZXJ0QXV0aDAgFw0yMTAzMDEyMzQ5MDVaGA8yMTIxMDMwMTIz
NDgxNVowEDEOMAwGA1UEAxMFQWxpY2UwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
ggEKAoIBAQC+TeMomPG1GyRev5Bqm6iV7D/csUkN2udPuXD/Ipp/f7MxwTk2F+G2
Pas4pDXPdspCfgNf+ZmItXo+G7pmGl96ZJIl4GdrHt1JSkkSj1tbGcSPgmZyfnK8
VBigudYZycndk1P83PaOWZvhhN0GoutuTooapxVFOakCUdm9HrOOnC6ICc9z3FZi
r4t6peXNXkChIFBpYuVpTpR13bmd9DfvUIMIpwDm/bNw8oS6u7MuAF2zVnHaQP8N
bDy28Klbsu/YeKgPJS7Lq6Yr/1hBCayR6K8LB092qhQ1uSjFZYeyiwLgkweMvBDV
ElPhPXxXe9mrlL8IpluLwydFXJMfuGFvAgMBAAGjcTBvMA4GA1UdDwEB/wQEAwID
uDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwHQYDVR0OBBYEFHtLPdX/
UG2dDrii6Tcj4lNN0/ewMB8GA1UdIwQYMBaAFB/986mJ3Of+uD0flHyHaMd0jp7U
MA0GCSqGSIb3DQEBCwUAA4ICAQBtCUi+lBOOxt6GrkiYcJxNXLYUmRPrb9UwsBkc
5fo7eZfG2ctZhY4S6VS41KFAthD/0OaA8Wqu1zreCI2c+pEnO+sdPvOWXk6LiyWo
J49pDVKl9dhN9WJb2GR03lPar62L1u19jZIe1RfyqoAbzcywsg0FQGM1hW59BZ1T
xaIlye9ulcasLbz+M9IKui7sVfYp7PI+shIbPGfQnp982EyS3X4gquWRCYzT9Qli
y1x4lbgv4B65YvAYdb2kzuJNBKUqAUwR/7HXzhrUih2+blJ20o89WWcKoJrZVGWL
jGwZlI0GX3TtIPpkuUZn1LgTqHcTW7yhtFPjtviLy2yHa4BJY7LbJE+uAZYHYO9Y
NQxdCztSQN/IPKPFnMdTTt8tsE/nxyn6Ukblb7Ko1sFSWfk6SgH++76t7fwg5qwe
aGhkTit91RmsGcc8i7lGRQO/tLMypnyxH8rIDXxD+6QTEOZAYsBy0QiM1pVK4Jie
xcvbw+Yz9O0eb15iFtPdOAVWOOVNtECabyOM7Wmd8XePxBvS8548DXBjzB2Ubnto
QRCqRyvSu4wa+VeWEsJMtDEyCUdDdd9dRBDKp9buW9hprIW5XvyGBrlKKmJnmz7k
yjS5PagH6T1bfmzJ7kdJ5OYk8N7WzJtDoqIntVF2D9J38JXacrClTK6nunftxlkF
dthKxA==
-----END CERTIFICATE-----`

	DATABASE_CLIENT_KEY = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAvk3jKJjxtRskXr+Qapuolew/3LFJDdrnT7lw/yKaf3+zMcE5
Nhfhtj2rOKQ1z3bKQn4DX/mZiLV6Phu6ZhpfemSSJeBnax7dSUpJEo9bWxnEj4Jm
cn5yvFQYoLnWGcnJ3ZNT/Nz2jlmb4YTdBqLrbk6KGqcVRTmpAlHZvR6zjpwuiAnP
c9xWYq+LeqXlzV5AoSBQaWLlaU6Udd25nfQ371CDCKcA5v2zcPKEuruzLgBds1Zx
2kD/DWw8tvCpW7Lv2HioDyUuy6umK/9YQQmskeivCwdPdqoUNbkoxWWHsosC4JMH
jLwQ1RJT4T18V3vZq5S/CKZbi8MnRVyTH7hhbwIDAQABAoIBABh+R9Vq0flAVA8J
0tmVzT32mUWbX86ztP/S21NLLd1pHzQxLV7j9f4Rs97na7GXFXM3atHIe1VYHjKu
OPB5Rn6nQRQ07Lqhz1Nmhz8nSlfQWjcqTmtAER5iKIVPRdot2Zh2JOIvwmAR8nk1
M4luIKUi4K2LgbZWNIWC6WZo9C1PffB56W4lxVJppPtNvtiop+S5cHnKmi27QUT2
f2DTFhjvC9fpJ255Cb3BwVk/Z1rkfSvL/de/eSeRjlxCADRsdZQ9NLZ5PNINNP4R
Y+8oIwhkDaYENIa0OpIGfCWzeafqQChO+OGYStfQgNYTY2uQjdyCQ5R7XW7ajR4b
hbmqrNECgYEA3Q82uGD5u6kfZBTfYrHUNmw45G6LQCseKQ8JT6BUSrxEae/Bxnec
B1AFbEEl+1mgPx7V5BkQ4ouzrmKnvMYaTqugfyZeiFGjTugBuCK2OhrDNRkUduWi
sArCELHSC01xbogXsHIxQYBsQT9tMawOgIPzpLDIAfjiz5A6yu1veukCgYEA3GI3
8fRJWrYx/z8Weqh31/df/sZGNeEvNJ0h6rl5t/ch7R8tipHYNkGokA2EVvl/eAMX
dqqZH9gLWinn2gcdScbvvjWXuH7jQI4+sfpUOwActD6tW+R/gdeXwqLmJdNdA66M
uZJEQ4lbDxBHpGKHIfJcmdHgS3wdUOtogIEukpcCgYEAg1OEeKjz8+6XPYfA5gsr
PWGxstORWn/DHTVXnLxtFzEdm6nZ/qQUR6vcbOGPRPGL57PT8fyKMWw5RMJLUDSI
cHA5mxAaXPXTBQ+D5faJN5+qlKLfq3rk2Zyqiex6EkjhuH6VRDey120J4wFhR38l
5md9mTuEttc7N4BBzUidT/ECgYEAwVIDQ2KIKmOijgY4YNaBclhUw/gHxOHI9/1S
sEWNSVwnTYtI9sIjCM0l4V+KFIV+VLdZkMXr1qw8oRYbhP0yqvIwggkfEz9zd8CP
rK4rzym1BEPq7K4PT8XgIWOmQc4cTMuENJDjAt9tmlQslD29zoB8zI33lB/G06H9
JKStRvcCgYEAznZ/Dh8D6lj5qKdQ3hzYMXgYTpwYhVshALWXkTVUy0f8gavtAR1R
Di1zoeYRMbflwbD+6/rWseLCVFpka0nRnUlAxt9i+dmZds3Jo2EfQnChNQZyvB3u
e5DVq/cK1VvLvDrxKZ1XGXgx/bPkT2QsXhbZaRexqsneorwaJ/EUwvs=
-----END RSA PRIVATE KEY-----`

	EXPIRED_CERT_FROM_ANOTHER_CA = `-----BEGIN CERTIFICATE-----
MIIEIzCCAgugAwIBAgIRAKRY04+EGcnAfPpeLf3dxZcwDQYJKoZIhvcNAQELBQAw
EjEQMA4GA1UEAxMHbXlzcWxDQTAeFw0xODEwMTYyMjU2NDdaFw0xODEwMTYyMzU2
NDdaMBUxEzARBgNVBAMTCmV4cGlyZVNvb24wggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQCZ/8fZc0q1I03L78hro9jr987Tn6hsJNGod61GiWOHybHezt5i
+SOp7S/fdmLQsRSopUOSmlAH6ta5QbffGbtHY2NJQKJq7N8KRt6aSfbHxPDG96Rp
Q0OZZLyiEaFz2jECoTjqwZX9duG5wA1/AVnZEKqnbAWdIWP9AOTzwdJ/ne4CLzyj
Lm/HUNi9xsZvU5xgb8ZSW3z8SOf39UedocmDcA/rTZWAkO6ELPvx4KD6t5aBC4ir
k7tGveQFxTvziZr3lNZk+NTX2OWUrz5yoH/nMiXtHe4JuytFsN5DYF1f6/3Fxl19
AhkCkxTj238/FFLID34W7mfZbgN59ByBgnPxAgMBAAGjcTBvMA4GA1UdDwEB/wQE
AwIDuDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwHQYDVR0OBBYEFJuI
c/+CukBgDTwmU6A+6++WwzaKMB8GA1UdIwQYMBaAFFmJZ4pg0d2UuDKpQzF3XQAM
LandMA0GCSqGSIb3DQEBCwUAA4ICAQAh7+THwf0fe3syAaPVqnpx2kswUAqP9VTw
waxXswwp632JnQa9vctuVBQ7DNwOHSixaNlM7yR+w1FlubwLzNRR5EXOgi2kl5Le
mewKBmJLpMwkmAbpCUB2B2ofJJguMe0JVQC6OC3eA3JsTc1/FtqJ4H1+RD5xT6hx
uOxla3zwfynYD4WdRMAosYVJouCScgWJpK+MWEkMCx94GUcO4Ik9acWhzBcdgaUG
qjbtTq5dHgVwernhJaiuUC2R5wEvb3rkhav2TYHJucFm0NHFbMCCYNbFAp1t1OyW
hiNrGtUGN2jBoFZ9OEZaWuY00mKs0Elp5/ugHQ5hW6HXam/4Fh95PMBR1QC+c5AC
AhdCYEXpZXkjCe5vnXHegBxAMV2FU33G9rPWWAi76sBlqjApGaYfbYJW63bhEOZT
AtnHlrPVw/GM16KkzMEEbi4lRvY4F3F2FJ+LZSMKMNs9aX/CAAWs9up3n7PcePP0
fV70C2hVtCJbIfRPaWvrVAAktBP9xLTnzUvzijPLMEJ9o45vWdrtvyBFknQCpMts
lw6sWU26m2gvxs6CcX3yt0bt8SxjqyulqrOdFCVSjZbGMDaIamdEKnC6k5ySyizn
SM2qNm+nV5FhjsyMyzs6OuCNEZGDAqklWBAHHqLncb6elO9NZgDysB/xn6jS+zqT
F1Y5M6wvLA==
-----END CERTIFICATE-----`
)

var _ = Describe("MySQLConnectionStringBuilder", func() {
	Describe("Build", func() {
		var (
			mysqlConnectionStringBuilder *db.MySQLConnectionStringBuilder
			mySQLAdapter                 *fakes.MySQLAdapter
			config                       db.Config
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

			mySQLAdapter = &fakes.MySQLAdapter{}

			mysqlConnectionStringBuilder = &db.MySQLConnectionStringBuilder{
				MySQLAdapter: mySQLAdapter,
			}
			mySQLAdapter.ParseDSNStub = func(dsn string) (cfg *mysql.Config, err error) {
				return mysql.ParseDSN(dsn)
			}
		})

		It("builds a connection string", func() {
			connectionString, err := mysqlConnectionStringBuilder.Build(config)
			Expect(err).NotTo(HaveOccurred())
			Expect(connectionString).To(Equal("some-user:some-password@tcp(some-host:1234)/some-database?parseTime=true&readTimeout=5s&timeout=5s&writeTimeout=5s"))
		})

		Context("when mysql.ParseDSN can't parse the connection string", func() {
			BeforeEach(func() {
				mySQLAdapter.ParseDSNReturns(nil, errors.New("foxtrot"))
			})

			It("returns an error", func() {
				_, err := mysqlConnectionStringBuilder.Build(config)
				Expect(err).To(MatchError("parsing db connection string: foxtrot"))
			})
		})

		Context("when requiring ssl", func() {
			var (
				caCertPool *x509.CertPool
			)

			BeforeEach(func() {
				caCertFile, err := ioutil.TempFile("", "")
				_, err = caCertFile.Write([]byte(DATABASE_CA_CERT))
				Expect(err).NotTo(HaveOccurred())

				config.RequireSSL = true
				config.CACert = caCertFile.Name()

				caCertPool = x509.NewCertPool()
				ok := caCertPool.AppendCertsFromPEM([]byte(DATABASE_CA_CERT))
				Expect(ok).To(BeTrue())
			})

			It("builds a tls connection string", func() {
				connectionString, err := mysqlConnectionStringBuilder.Build(config)
				Expect(err).NotTo(HaveOccurred())
				Expect(connectionString).To(Equal("some-user:some-password@tcp(some-host:1234)/some-database?parseTime=true&readTimeout=5s&timeout=5s&tls=some-database-tls&writeTimeout=5s"))

				Expect(mySQLAdapter.RegisterTLSConfigCallCount()).To(Equal(1))
				passedTLSConfigName, passedTLSConfig := mySQLAdapter.RegisterTLSConfigArgsForCall(0)
				Expect(passedTLSConfigName).To(Equal("some-database-tls"))
				Expect(passedTLSConfig).To(Equal(&tls.Config{
					InsecureSkipVerify: false,
					RootCAs:            caCertPool,
				}))
			})

			Context("when SkipHostnameValidation is true", func() {
				BeforeEach(func() {
					config.SkipHostnameValidation = true
				})

				It("builds tls config skipping hostname", func() {
					connectionString, err := mysqlConnectionStringBuilder.Build(config)
					Expect(err).NotTo(HaveOccurred())
					Expect(connectionString).To(Equal("some-user:some-password@tcp(some-host:1234)/some-database?parseTime=true&readTimeout=5s&timeout=5s&tls=some-database-tls&writeTimeout=5s"))

					Expect(mySQLAdapter.RegisterTLSConfigCallCount()).To(Equal(1))
					passedTLSConfigName, passedTLSConfig := mySQLAdapter.RegisterTLSConfigArgsForCall(0)
					Expect(passedTLSConfigName).To(Equal("some-database-tls"))
					Expect(passedTLSConfig.InsecureSkipVerify).To(BeTrue())
					Expect(passedTLSConfig.RootCAs).To(Equal(caCertPool))
					Expect(passedTLSConfig.Certificates).To(BeNil())
					// impossible to assert VerifyPeerCertificate is set to a specfic function
					Expect(passedTLSConfig.VerifyPeerCertificate).NotTo(BeNil())
				})
			})

			Context("when it can't read the ca cert file", func() {
				BeforeEach(func() {
					config.CACert = "/foo/bar"
				})

				It("returns an error", func() {
					_, err := mysqlConnectionStringBuilder.Build(config)
					Expect(err).To(MatchError("reading db ca cert file: open /foo/bar: no such file or directory"))
				})
			})

			Context("when it can't append the ca cert to the cert pool", func() {
				BeforeEach(func() {
					caCertFile, err := ioutil.TempFile("", "")
					_, err = caCertFile.Write([]byte("bad cert"))
					Expect(err).NotTo(HaveOccurred())

					config.CACert = caCertFile.Name()
				})

				It("returns an error", func() {
					_, err := mysqlConnectionStringBuilder.Build(config)
					Expect(err).To(MatchError("appending cert to pool from pem - invalid cert bytes"))
				})
			})

			Context("when it can't register TLS config", func() {
				BeforeEach(func() {
					mySQLAdapter.RegisterTLSConfigReturns(errors.New("bad things happened"))
				})

				It("retruns an error", func() {
					_, err := mysqlConnectionStringBuilder.Build(config)
					Expect(err).To(MatchError("registering mysql tls config: bad things happened"))
				})
			})
		})
	})

	Describe("VerifyCertificatesIgnoreHostname", func() {
		var (
			caCertPool *x509.CertPool
		)

		BeforeEach(func() {
			caCertPool = x509.NewCertPool()
			ok := caCertPool.AppendCertsFromPEM([]byte(DATABASE_CA_CERT))
			Expect(ok).To(BeTrue())
		})

		It("verifies that provided certificates are valid", func() {
			block, _ := pem.Decode([]byte(DATABASE_CLIENT_CERT))

			err := db.VerifyCertificatesIgnoreHostname([][]byte{
				block.Bytes,
			}, caCertPool)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when raw certs are not parsable", func() {
			It("returns an error", func() {
				err := db.VerifyCertificatesIgnoreHostname([][]byte{
					[]byte("foo"),
					[]byte("bar"),
				}, nil)
				Expect(err.Error()).To(ContainSubstring("tls: failed to parse certificate from server: asn1: structure error: tags don't match"))
			})
		})

		Context("when verifying an expired cert", func() {
			It("returns an error", func() {
				block, _ := pem.Decode([]byte(EXPIRED_CERT_FROM_ANOTHER_CA))

				err := db.VerifyCertificatesIgnoreHostname([][]byte{
					block.Bytes,
				}, caCertPool)

				Expect(err.Error()).To(ContainSubstring("x509: certificate has expired or is not yet valid"))
			})
		})
	})
})

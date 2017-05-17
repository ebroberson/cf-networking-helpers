package timeouts_test

import (
	"context"
	"fmt"
	"math/rand"
	"os/exec"
	"policy-server/integration/helpers"
	"time"

	"code.cloudfoundry.org/go-db-helpers/db"
	"code.cloudfoundry.org/go-db-helpers/testsupport"

	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var createTable = `CREATE TABLE IF NOT EXISTS mytable ( id SERIAL PRIMARY KEY);`
var testTimeoutInSeconds = float64(5)

var _ = Describe("Timeout", func() {
	var (
		testDatabase     *testsupport.TestDatabase
		dbName           string
		dbConnectionInfo *testsupport.DBConnectionInfo
		ctx              context.Context
	)

	BeforeEach(func() {
		dbName = fmt.Sprintf("test_%x", rand.Int())
		dbConnectionInfo = testsupport.GetDBConnectionInfo()
	})

	AfterEach(func() {
		if testDatabase != nil {
			testDatabase.Destroy()
			testDatabase = nil
		}
	})

	var database *sqlx.DB
	Describe("postgres", func() {
		dbConnectionInfo = testsupport.GetDBConnectionInfo()
		if dbConnectionInfo.Type != "postgres" {
			Skip("skipping postgres tests")
		}

		Context("when the read timeout is greater than the context timeout and the database is unreachable", func() {
			BeforeEach(func() {
				ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
				// dbConnectionInfo.ConnectTimeout = 3 * time.Second // not needed for QueryRowContext
				dbConnectionInfo.ReadTimeout = 3 * time.Second
				// dbConnectionInfo.WriteTimeout = 3 * time.Second // not needed for QueryRowContext
				testDatabase = dbConnectionInfo.CreateDatabase(dbName)

				var err error
				database, err = db.GetConnectionPool(testDatabase.DBConfig())
				Expect(err).NotTo(HaveOccurred())

				By("creating a table")
				_, err = database.Exec(createTable)
				Expect(err).NotTo(HaveOccurred())

				By("blocking access to port " + dbConnectionInfo.Port)
				mustSucceed("iptables", "-A", "INPUT", "-p", "tcp", "--dport", dbConnectionInfo.Port, "-j", "DROP")
			})

			AfterEach(func() {
				By("allowing access to port " + dbConnectionInfo.Port)
				mustSucceed("iptables", "-D", "INPUT", "-p", "tcp", "--dport", dbConnectionInfo.Port, "-j", "DROP")
			})

			Describe("QueryRowContext", func() {
				It("returns a context deadline exceeded error", func(done Done) {
					defer database.Close()
					var databaseName string
					err := database.QueryRowContext(ctx, "SELECT current_database();").Scan(&databaseName)
					Expect(err).To(HaveOccurred())
					Expect(err).To(BeAssignableToTypeOf(context.DeadlineExceeded))

					close(done)
				}, testTimeoutInSeconds)
			})

			Describe("QueryContext", func() {
				It("returns a context deadline exceeded error", func(done Done) {
					defer database.Close()
					_, err := database.QueryContext(ctx, "SELECT id FROM mytable;")
					Expect(err).To(HaveOccurred())
					Expect(err).To(BeAssignableToTypeOf(context.DeadlineExceeded))

					close(done)
				}, testTimeoutInSeconds)
			})

			Describe("ExecContext", func() {
				PIt("returns a context deadline exceeded error", func(done Done) {
					defer database.Close()

					_, err := database.Exec("INSERT into mytable (id) values (1);")
					// _, err := database.ExecContext(ctx, "INSERT into mytable (id) values (1);")
					Expect(err).NotTo(HaveOccurred())
					// Expect(err).To(BeAssignableToTypeOf(context.DeadlineExceeded))

					close(done)
				}, testTimeoutInSeconds)
			})

		})
	})
})

func mustSucceed(binary string, args ...string) string {
	cmd := exec.Command(binary, args...)
	sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, helpers.DEFAULT_TIMEOUT).Should(gexec.Exit(0))
	return string(sess.Out.Contents())
}

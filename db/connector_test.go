package db_test

import (
	"fmt"
	"math/rand"

	"code.cloudfoundry.org/cf-networking-helpers/db"
	"code.cloudfoundry.org/cf-networking-helpers/testsupport"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetConnectionPool", func() {
	var (
		dbConf db.Config
	)

	BeforeEach(func() {
		dbConf = testsupport.GetDBConfig()
		dbConf.DatabaseName = fmt.Sprintf("test_%x", rand.Int())
		testsupport.CreateDatabase(dbConf)
	})

	AfterEach(func() {
		dbConf = testsupport.GetDBConfig()
		testsupport.RemoveDatabase(dbConf)
	})

	It("returns an error if the connection string cannot be created", func() {
		_, err := db.GetConnectionPool(db.Config{})
		Expect(err).To(MatchError("failed to create connection string: timeout must be at least 1 second: 0"))
	})

	It("returns a database reference", func() {
		database, err := db.GetConnectionPool(dbConf)
		Expect(err).NotTo(HaveOccurred())
		defer database.Close()

		var databaseName string
		if database.DriverName() == "postgres" {
			err = database.QueryRow("SELECT current_database();").Scan(&databaseName)
		} else if database.DriverName() == "mysql" {
			err = database.QueryRow("SELECT DATABASE();").Scan(&databaseName)
		} else {
			panic("unsupported db type")
		}
		Expect(err).NotTo(HaveOccurred())
		Expect(databaseName).To(Equal(dbConf.DatabaseName))
	})

	Context("when the database cannot be accessed", func() {
		It("returns a non-retryable error", func() {
			testsupport.RemoveDatabase(dbConf)
			_, err := db.GetConnectionPool(dbConf)
			Expect(err).To(HaveOccurred())

			Expect(err).NotTo(BeAssignableToTypeOf(db.RetriableError{}))
			Expect(err).To(MatchError(ContainSubstring("unable to ping")))
		})
	})

	Context("when there is a network connectivity problem", func() {
		It("returns a retriable error", func() {
			dbConf.Port = 0

			_, err := db.GetConnectionPool(dbConf)
			Expect(err).To(HaveOccurred())

			Expect(err).To(BeAssignableToTypeOf(db.RetriableError{}))
			Expect(err.Error()).To(ContainSubstring("unable to ping"))
		})
	})

	It("sets the databaseConfig.Type as the DriverName", func() {
		database, err := db.GetConnectionPool(dbConf)
		Expect(err).NotTo(HaveOccurred())
		defer database.Close()

		Expect(database.DriverName()).To(Equal(dbConf.Type))
	})
})

package testsupport

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"code.cloudfoundry.org/cf-networking-helpers/db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func CreateDatabase(config db.Config) {
	_, err := execSQL(config, fmt.Sprintf("CREATE DATABASE %s", config.DatabaseName))
	Expect(err).NotTo(HaveOccurred())
}

func RemoveDatabase(config db.Config) error {
	_, err := execSQL(config, fmt.Sprintf("DROP DATABASE %s", config.DatabaseName))
	return err
}

func execSQL(c db.Config, sqlCommand string) (string, error) {
	var cmd *exec.Cmd

	if c.Type == "mysql" {
		cmd = exec.Command("mysql",
			"-h", c.Host,
			"-P", strconv.Itoa(int(c.Port)),
			"-u", c.User,
			"-e", sqlCommand)
		cmd.Env = append(os.Environ(), "MYSQL_PWD="+c.Password)
	} else if c.Type == "postgres" {
		cmd = exec.Command("psql",
			"-h", c.Host,
			"-p", strconv.Itoa(int(c.Port)),
			"-U", c.User,
			"-c", sqlCommand)
		cmd.Env = append(os.Environ(), "PGPASSWORD="+c.Password)
	} else {
		panic("unsupported database type: " + c.Type)
	}

	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, "9s").Should(gexec.Exit())
	if session.ExitCode() != 0 {
		return "", fmt.Errorf("unexpected exit code: %d", session.ExitCode())
	}
	return string(session.Out.Contents()), nil
}

const DefaultDBTimeout = 5

func getPostgresDBConfig() db.Config {
	return db.Config{
		Type:     "postgres",
		User:     "postgres",
		Password: "",
		Host:     "127.0.0.1",
		Port:     5432,
		Timeout:  DefaultDBTimeout,
	}
}

func getMySQLDBConfig() db.Config {
	return db.Config{
		Type:     "mysql",
		User:     "root",
		Password: "password",
		Host:     "127.0.0.1",
		Port:     3306,
		Timeout:  DefaultDBTimeout,
	}
}

func GetDBConfig() db.Config {
	switch os.Getenv("DB") {
	case "mysql":
		return getMySQLDBConfig()
	case "postgres":
		return getPostgresDBConfig()
	default:
		panic("unable to determine database to use.  Set environment variable DB")
	}
}

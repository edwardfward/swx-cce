package cce

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

type appData struct {
	db      *sql.DB
	connStr string
}

// TODO: Add more robust data and statistics
func (d *appData) dbInfo() string {
	return d.connStr
}

// TODO: Add more robust error checking and default env vars
func createConnStr() (string, error) {

	vars := make(map[string]string)

	connStr := ""

	vars["dbname"] = os.Getenv("DB_NAME")
	vars["user"] = os.Getenv("DB_USERNAME")
	vars["password"] = os.Getenv("DB_PASSWORD")
	vars["host"] = os.Getenv("DB_HOST")
	vars["port"] = os.Getenv("DB_PORT")
	vars["sslmode"] = os.Getenv("DB_SSL_MODE")
	vars["connect_timeout"] = os.Getenv("DB_CONNECT_TIMEOUT")
	vars["sslcert"] = os.Getenv("DB_SSL_CERT")
	vars["sslkey"] = os.Getenv("DB_SSL_KEY")
	vars["sslrootcert"] = os.Getenv("DB_SSL_ROOT_CERT")

	for k, v := range vars {
		if v == "" {
			continue
		} else {
			connStr += fmt.Sprintf("%s=%s ", k, v)
		}
	}

	return strings.TrimSpace(connStr), nil

}

func (d *appData) testDatabase() string {
	result, err := d.db.Exec(`SELECT $1;`, 1)
	if err != nil {
		return fmt.Sprintf("ERROR: Database is not functioning properly (%v)",
			err)
	} else {
		return fmt.Sprintf("RESULT: %v", result)
	}
}

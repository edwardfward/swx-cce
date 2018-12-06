package app

import (
	"database/sql"
	"fmt"
	"os"
)

type appData struct {
	db      *sql.DB
	connStr string
}

// TODO: Add more robust data and statistics
func (s *server) dbInfo() string {
	return s.data.connStr
}

// TODO: Add more robust error checking and default env vars
func CreateConnStr() (string, error) {
	// gather env vars for postgres database
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dbUsername, dbPassword, dbHost, dbPort, dbName), nil
}

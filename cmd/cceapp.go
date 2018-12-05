package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var (
	db         *sql.DB
	dbHost     = os.Getenv("DB_HOST")
	dbName     = os.Getenv("DB_NAME")
	dbPort     = os.Getenv("DB_PORT")
	dbUsername = os.Getenv("DB_USERNAME")
	dbPassword = os.Getenv("DB_PASSWORD")
)

func main() {

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dbUsername, dbPassword, dbHost, dbPort, dbName)

	log.Printf("Connection string: %s", connStr)

	var err error

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("DATABASE CONNECTION ERROR: Could not connect to database"+
			" %s on %s:%s", dbName, dbHost, dbPort)
	}

	defer db.Close()

	http.HandleFunc("/info", displayEnvVars)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("SERVER ERROR: Server unexpectedly quit")
	}
}

func displayEnvVars(w http.ResponseWriter, r *http.Request) {

	log.Printf("/info request recieved")

	dbStats := db.Stats()

	n, err := fmt.Fprintf(w, "Database Url: %s\n"+
		"Database Port: %s\n"+
		"Database Username: %s\n"+
		"Database Password: %s\n"+
		"Database Max Connections: %d\n", dbHost, dbPort, dbUsername, dbPassword,
		dbStats.OpenConnections)

	if n == 0 || err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("500 - Database parameters did not write properly"))
		if err != nil {
			return
		}
	}
}

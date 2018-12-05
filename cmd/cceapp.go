package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("/info", displayEnvVars)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func displayEnvVars(w http.ResponseWriter, r *http.Request) {
	dbUrl := os.Getenv("DB_URL")
	dbPort := os.Getenv("DB_PORT")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	n, err := fmt.Fprintf(w, "Database Url: %s\n"+
		"Database Port: %s\n"+
		"Database Username: %s\n"+
		"Database Password: %s\n", dbUrl, dbPort, dbUsername, dbPassword)

	if n == 0 || err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("500 - Database parameters did not write properly"))
		if err != nil {
			return
		}
	}
}

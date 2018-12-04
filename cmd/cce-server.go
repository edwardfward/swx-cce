package main

import (
	"fmt"
	"os"
)

func main(){

	dbUrl := os.Getenv("DB_URL")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")

	fmt.Printf("Database Url: %s\n" +
		"Database Username: %s\n" +
		"Database Password: %s\n" +
		"Database Port: %s\n", dbUrl, dbUsername, dbPassword, dbPort)
}

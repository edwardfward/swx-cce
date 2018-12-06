package cce

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

type server struct {
	data   *appData
	router *http.ServeMux
}

func NewServer() (*server, error) {

	// delay to start to ensure Postgres container is ready
	time.Sleep(5 * time.Second)

	// build connection string
	// TODO: incorporate SSL options
	connStr, err := createConnStr()
	if err != nil {
		log.Fatalf("ERROR: Connection string failed (%s)", connStr)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("ERROR: Database failed to connect (%s)", err)
	}

	s := server{}
	s.data = &appData{}
	s.data.db = db
	s.data.connStr = connStr
	s.router = http.NewServeMux()

	return &s, nil
}

// Close closes the connection to the database.
func (s *server) Close() {
	err := s.data.db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// Start server
func (s *server) StartServer() {
	s.addRouting()
	err := http.ListenAndServe(":8080", s.router)
	if err != nil {
		log.Fatalf("ERROR: Failed to start (%v)", err)
	}
	log.Print("SERVER: Started and listening on port :8080")

}

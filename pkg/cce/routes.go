package cce

import (
	"fmt"
	"log"
	"net/http"
)

func (s *server) addRouting() {
	s.router.HandleFunc("/", s.handleIndex)
	s.router.HandleFunc("/info", s.handleInfo)
	s.router.HandleFunc("/testdb", s.handleTestDb)
}

func (s *server) handleIndex(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		n, err := fmt.Fprint(w, "Welcome to the CCE Input Prototype")
		if n == 0 || err != nil {
			log.Print(err)
		}
	case "POST":
		http.Error(w, "Unauthorized method", http.StatusMethodNotAllowed)
	}
}

func (s *server) handleInfo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		n, err := fmt.Fprint(w, s.dbInfo())
		if n == 0 || err != nil {
			http.Error(w, "Unable to write database info",
				http.StatusInternalServerError)
		}
	case "POST":
		http.Error(w, "Unauthorized method", http.StatusMethodNotAllowed)
	}
}

func (s *server) handleTestDb(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		n, err := fmt.Fprint(w, s.testDatabase())
		if n == 0 || err != nil {
			http.Error(w, "Unable to read from database",
				http.StatusInternalServerError)
		}
	case "POST":
		http.Error(w, "Unauthorized method", http.StatusMethodNotAllowed)
	}
}

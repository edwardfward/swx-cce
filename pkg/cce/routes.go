package cce

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (s *server) addRouting() {
	s.router.HandleFunc("/", s.handleIndexPage)
	s.router.HandleFunc("/admin", s.handleAdminPage)
	s.router.HandleFunc("/cce/{cce}", s.handleCCEPage)

	s.router.HandleFunc("/api/v1/{cce}[a-zA-Z]+", s.cceHandler)
	s.router.HandleFunc("/api/v1/{cce}/{cceLimitation:[0-9]+}",
		s.cceLimitHandler)
}

func (s *server) handleIndexPage(w http.ResponseWriter, r *http.Request) {
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

func (s *server) handleAdminPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		n, err := fmt.Fprint(w, "ADMIN PAGE")
		if n == 0 || err != nil {
			http.Error(w, "Unable to access admin page",
				http.StatusInternalServerError)
		}
	case "POST":
		http.Error(w, "Unauthorized method", http.StatusMethodNotAllowed)
	}
}

func (s *server) handleCCEPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	switch r.Method {
	case "GET":
		n, err := fmt.Fprintf(w, "%s", vars["cce"])
		if n == 0 || err != nil {
			http.Error(w, "Unable to read cce variable",
				http.StatusInternalServerError)
		}
	case "POST":
		http.Error(w, "Unauthorized method", http.StatusMethodNotAllowed)
	}
}

// TODO: flush out handlerfunc
func (s *server) cceHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	switch r.Method {
	case "GET":
		n, err := fmt.Fprintf(w, "CCE: %s", vars["cce"])
		if n == 0 || err != nil {
			http.Error(w, "Unable to parse variable from url",
				http.StatusInternalServerError)
		}
	}
}

// TODO: flush out handlerfunc
func (s *server) cceLimitHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	switch r.Method {
	case "GET":
		n, err := fmt.Fprintf(w, "CCE: %s Limit Id: %s",
			vars["cce"], vars["cceLimitation"])
		if n == 0 || err != nil {
			http.Error(w, "Unable to parse variables from url",
				http.StatusInternalServerError)
		}
	}
}

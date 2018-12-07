package cce

import (
	"net/http"
)

func (s *server) addRouting() {
	// web routing
	s.router.HandleFunc("/", s.handleIndex)
	s.router.HandleFunc("/admin", s.handleAdmin)
	s.router.HandleFunc("/cce/{cce}", s.handleCCE)

	// static routing
	fs := http.FileServer(http.Dir("./web/static"))
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// api routing
	s.router.HandleFunc("/api/v1/{cce:[a-zA-Z]+}", s.apiManageCCEs)
	s.router.HandleFunc("/api/v1/{cce}/{cceLimitation:[0-9]+}",
		s.apiManageLimits)
}

func (s *server) handleIndex(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
	case "POST":
	}
}

func (s *server) handleAdmin(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
	case "POST":
	}
}

func (s *server) handleCCE(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
	case "POST":
	}
}

// apiManageCCEs allows facilitator or admin to perform ACID operations on
// a CCE. Function queries database and returns CCE JSON.
// TODO: build test and flush out
func (s *server) apiManageCCEs(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
	case "POST":

	}
}

// apiManageLimits allows users, facilitators, and admins to perform ACID
// operations on individual CCE limits. Function returns and accepts
// TODO: flush out
func (s *server) apiManageLimits(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
	case "POST":
	}
}

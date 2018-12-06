package cce

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *server) addRouting() {
	// web routing
	s.router.HandleFunc("/", s.handleIndexPage)
	s.router.HandleFunc("/admin", s.handleAdminPage)
	s.router.HandleFunc("/cce/{cce}", s.handleCCEPage)

	// static routing
	fs := http.FileServer(http.Dir("./web/static"))
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// api routing
	s.router.HandleFunc("/api/v1/{cce:[a-zA-Z]+}", s.apiCCEHandler)
	s.router.HandleFunc("/api/v1/{cce}/{cceLimitation:[0-9]+}",
		s.apiCCELimitHandler)
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
	log.Print("CCE: Accessing page...")

	cce := &CCE{Title: vars["cce"]}
	for i := range [100]int{} {
		cce.Limits = append(cce.Limits,
			Limitation{Limit: "System is not secure enough.",
				Id: i + 1})

	}

	tmpl, err := template.ParseFiles("web/templates/cce.html")
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	switch r.Method {
	case "GET":
		err := tmpl.Execute(w, cce)
		if err != nil {
			// http.Error(w, "Unable to read cce variable",
			//	http.StatusInternalServerError)
		}
	case "POST":
		http.Error(w, "Unauthorized method", http.StatusMethodNotAllowed)
	}
}

// TODO: flush out
func (s *server) apiCCEHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		n, err := fmt.Fprint(w, "Test")
		if n == 0 || err != nil {
			http.Error(w, "Unable to parse template",
				http.StatusInternalServerError)
		}

	case "POST":

	}
}

// TODO: flush out
func (s *server) apiCCELimitHandler(w http.ResponseWriter, r *http.Request) {

	cce := &CCE{}
	testJson, err := json.Marshal(cce)
	if err != nil {
		log.Print("ERROR: CCELimitation JSON did not marshal")
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		n, err := w.Write(testJson)
		if n == 0 || err != nil {
			http.Error(w, "Unable to parse variables from url",
				http.StatusInternalServerError)
		}
	}
}

package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type app struct {
	config  config
	repo    *repo
	handler *mux.Router
}

func build(c config) (*app, error) {
	a := &app{config: c}

	r, err := buildRepo(c)
	if err != nil {
		return nil, err
	}
	a.repo = r

	m := mux.NewRouter()
	m.Methods("GET").Path("/api/notes").HandlerFunc(a.GetNotes)
	m.Methods("POST").Path("/api/notes").HandlerFunc(a.PostNotes)
	m.Methods("GET").Path("/api/notes/{id}").HandlerFunc(a.GetNote)
	m.Methods("PUT").Path("/api/notes/{id}").HandlerFunc(a.PutNote)
	m.Methods("DELETE").Path("/api/notes/{id}").HandlerFunc(a.DeleteNote)
	a.handler = m

	return a, nil
}

func (a app) run() {
	var url string
	if strings.HasPrefix(a.config.address, ":") {
		url = "http://localhost" + a.config.address
	} else {
		url = "http://" + a.config.address
	}

	log.Println("Launching server on " + url)
	log.Fatal(http.ListenAndServe(a.config.address, a.handler))
}

func (a app) GetNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := a.repo.ListNotes()
	if err != nil {
		nok(w, 500, err)
		return
	}

	ok(w, notes)
	return
}

func (a app) PostNotes(w http.ResponseWriter, r *http.Request) {
	var note note
	if !in(w, r, &note) {
		return
	}

	id, err := a.repo.AddNote(&note)
	if err != nil {
		nok(w, 500, err)
		return
	}

	// construct output
	var output struct {
		ID int64 `json:"id"`
	}
	output.ID = id

	ok(w, output)
}

func (a app) GetNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	note, err := a.repo.ReadNote(id)
	if err != nil {
		nok(w, 500, err)
		return
	}

	ok(w, note)
}

func (a app) PutNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var note note
	if !in(w, r, &note) {
		return
	}

	err := a.repo.EditNote(id, &note)
	if err != nil {
		nok(w, 500, err)
		return
	}
	ok(w, nil)
}

func (a app) DeleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := a.repo.DeleteNote(id)
	if err != nil {
		nok(w, 500, err)
		return
	}
	ok(w, nil)
	return
}

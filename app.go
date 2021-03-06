package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type app struct {
	config  config
	repo    *repo
	handler http.Handler
}

func build(c config) (*app, error) {
	a := &app{config: c}

	r, err := buildRepo(c)
	if err != nil {
		return nil, err
	}
	a.repo = r

	m := mux.NewRouter()
	m.Methods("GET").Path(c.root + "/notes").HandlerFunc(a.GetNotes)
	m.Methods("POST").Path(c.root + "/notes").HandlerFunc(a.PostNotes)
	m.Methods("GET").Path(c.root + "/notes/{id}").HandlerFunc(a.GetNote)
	m.Methods("PUT").Path(c.root + "/notes/{id}").HandlerFunc(a.PutNote)
	m.Methods("DELETE").Path(c.root + "/notes/{id}").HandlerFunc(a.DeleteNote)
	a.handler = m

	if c.dev {
		a.handler = debug(a.handler)
	}

	return a, nil
}

func (a app) run() {
	log.Fatal(http.ListenAndServe(a.config.address, a.handler))
}

func debug(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
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

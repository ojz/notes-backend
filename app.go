package main

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

type app struct {
	config  config
	repo    *repo
	handler *http.ServeMux
}

func build(c config) (*app, error) {
	a := &app{config: c}

	r, err := buildRepo(c)
	if err != nil {
		return nil, err
	}
	a.repo = r

	m := http.NewServeMux()
	m.HandleFunc("/api/notes", a.notes)
	m.HandleFunc("/api/notes/", a.note)
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

func (a app) notes(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		notes, err := a.repo.ListNotes()
		if err != nil {
			nok(w, 500, err)
			return
		}

		ok(w, notes)
		return
	}

	if r.Method == "POST" {
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
}

func (a app) note(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[3]

	if r.Method == "GET" {
		note, err := a.repo.ReadNote(id)
		if err != nil {
			nok(w, 500, err)
			return
		}

		ok(w, note)
		return
	}

	// PUT
	if r.Method == "PUT" {
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
		return
	}

	// DELETE
	if r.Method == "DELETE" {
		err := a.repo.DeleteNote(id)
		if err != nil {
			nok(w, 500, err)
			return
		}
		ok(w, nil)
		return
	}

	nok(w, 500, errors.New("Method unsupported atm"))
}

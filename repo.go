package main

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type note struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body,omitempty"`
}

type repo struct {
	db *sqlx.DB

	listNotes,
	addNote, readNote, editNote, deleteNote *sqlx.Stmt
}

func buildRepo(c config) (*repo, error) {
	db := sqlx.MustConnect("sqlite3", c.database)

	if c.init {
		_, err := db.Exec(`CREATE TABLE notes (id INTEGER PRIMARY KEY, title TINYTEXT, body TINYTEXT)`)
		if err != nil {
			return nil, err
		}
	}

	prepare := func(sql string) *sqlx.Stmt {
		stmt, err := db.Preparex(sql)
		if err != nil {
			panic(err)
		}
		return stmt
	}

	r := repo{
		listNotes:  prepare("SELECT id, title FROM notes"),
		addNote:    prepare("INSERT INTO notes (title, body) VALUES (?, ?)"),
		readNote:   prepare("SELECT title, body FROM notes WHERE id = ?"),
		editNote:   prepare("UPDATE notes SET title = ?, body = ? WHERE id = ?"),
		deleteNote: prepare("DELETE FROM notes WHERE id = ?"),
	}

	return &r, nil
}

// add note

func (r repo) AddNote(note *note) (int64, error) {
	res, err := r.addNote.Exec(note.Title, note.Body)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// read note

func (r repo) ReadNote(id string) (*note, error) {
	var note note
	err := r.readNote.Get(&note, id)
	if err == sql.ErrNoRows {
		return nil, errors.New("note not found")
	}
	note.ID = id

	return &note, nil
}

// update note
func (r repo) EditNote(id string, note *note) error {
	_, err := r.editNote.Exec(note.Title, note.Body, id)
	return err
}

// delete note
func (r repo) DeleteNote(id string) error {
	_, err := r.deleteNote.Exec(id)
	return err
}

// list notes
func (r repo) ListNotes() ([]*note, error) {
	notes := []*note{}
	err := r.listNotes.Select(&notes)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

package domain

import (
	"errors"
	"time"
)

var ErrCreateNote = errors.New("failed to create a note")

type Service interface {
	Create(note *Note) (NoteID, error)
	GetOne(id int) (*Note, error)
	GetList(limit, offset int) ([]Note, error)
	Delete(id int) error
}

type Repository interface {
	Create(note *Note) (NoteID, error)
	GetOne(id int) (*Note, error)
	GetList(limit, offset int) ([]Note, error)
	Delete(id int) error
}

type NoteID int

type Note struct {
	ID        int
	Name      string
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

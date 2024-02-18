package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/synthao/notes/internal/domain"
	"time"
)

type oneDTO struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Text      string    `db:"text"`
	CreatedAt time.Time `db:"created_at"`
}

type listDTO struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) domain.Repository {
	return &repository{db: db}
}

func (r *repository) Create(note *domain.Note) (domain.NoteID, error) {
	args := map[string]interface{}{
		"name": note.Name,
		"text": note.Text,
	}

	exec, err := r.db.NamedExec("INSERT INTO notes(name, text) VALUES (:name, :text)", args)
	if err != nil {
		return 0, fmt.Errorf("%w, named exec, %w", domain.ErrCreateNote, err)
	}

	id, err := exec.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%w, get last insert id, %w", domain.ErrCreateNote, err)
	}

	return domain.NoteID(id), nil
}

func (r *repository) GetOne(id int) (*domain.Note, error) {
	var dest oneDTO

	err := r.db.Get(&dest, "SELECT id, name, text, created_at FROM notes WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	return &domain.Note{
		ID:        dest.ID,
		Name:      dest.Name,
		Text:      dest.Text,
		CreatedAt: dest.CreatedAt,
	}, nil
}

func (r *repository) GetList(limit, offset int) ([]domain.Note, error) {
	var dest []listDTO

	err := r.db.Select(&dest, "SELECT id, name, text FROM notes")
	if err != nil {
		return nil, err
	}

	return fromListDTOToDomain(dest), nil
}

func (r *repository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM notes WHERE id=?", id)
	if err != nil {
		return err
	}

	return nil
}

func fromListDTOToDomain(dto []listDTO) []domain.Note {
	res := make([]domain.Note, len(dto))

	for i, item := range dto {
		res[i] = domain.Note{
			ID:   item.ID,
			Name: item.Name,
		}
	}

	return res
}

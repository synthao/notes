package service

import (
	"github.com/synthao/notes/internal/domain"
	"go.uber.org/zap"
)

type service struct {
	logger     *zap.Logger
	repository domain.Repository
}

func NewService(logger *zap.Logger, repository domain.Repository) domain.Service {
	return &service{logger: logger, repository: repository}
}

func (s *service) Create(note *domain.Note) (domain.NoteID, error) {
	id, err := s.repository.Create(note)
	if err != nil {
		s.logger.Error(err.Error(), zap.Any("payload", note))
	}

	return id, nil
}

func (s *service) GetOne(id int) (*domain.Note, error) {
	return s.repository.GetOne(id)
}

func (s *service) GetList(limit, offset int) ([]domain.Note, error) {
	return s.repository.GetList(limit, offset)
}

func (s *service) Delete(id int) error {
	if _, err := s.repository.GetOne(id); err != nil {
		return err
	}

	if err := s.repository.Delete(id); err != nil {
		s.logger.Error(err.Error(), zap.Int("id", id))
		return err
	}

	return nil
}

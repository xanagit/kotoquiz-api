package services

import (
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/repositories"
)

type WordService interface {
	ReadWord(id string) (*models.Word, error)
	CreateWord(word *models.Word) error
	UpdateWord(word *models.Word) error
	DeleteWord(id string) error
}

type WordServiceImpl struct {
	Repo repositories.WordRepository
}

func (s *WordServiceImpl) ReadWord(id string) (*models.Word, error) {
	return s.Repo.ReadWord(id)
}

func (s *WordServiceImpl) CreateWord(word *models.Word) error {
	return s.Repo.CreateWord(word)
}

func (s *WordServiceImpl) UpdateWord(word *models.Word) error {
	return s.Repo.UpdateWord(word)
}

func (s *WordServiceImpl) DeleteWord(id string) error {
	return s.Repo.DeleteWord(id)
}

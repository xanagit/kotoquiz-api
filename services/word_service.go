package services

import (
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/repositories"
)

type WordService interface {
	ReadWord(id uuid.UUID) (*models.Word, error)
	CreateWord(word *models.Word) error
	UpdateWord(word *models.Word) error
	DeleteWord(id uuid.UUID) error
}

type WordServiceImpl struct {
	Repo repositories.WordRepository
}

// Make sure that WordServiceImpl implements WordService
var _ WordService = (*WordServiceImpl)(nil)

func (s *WordServiceImpl) ReadWord(id uuid.UUID) (*models.Word, error) {
	return s.Repo.ReadWord(id)
}

func (s *WordServiceImpl) CreateWord(word *models.Word) error {
	word.ID = uuid.Nil
	word.Translation.Type = models.Translation
	for _, t := range word.Tags {
		t.Type = models.Tag
	}
	for _, t := range word.Levels {
		t.Category.Type = models.Category
		for _, l := range t.LevelNames {
			l.Type = models.LevelName
		}
	}
	return s.Repo.CreateWord(word)
}

func (s *WordServiceImpl) UpdateWord(word *models.Word) error {
	return s.Repo.UpdateWord(word)
}

func (s *WordServiceImpl) DeleteWord(id uuid.UUID) error {
	return s.Repo.DeleteWord(id)
}

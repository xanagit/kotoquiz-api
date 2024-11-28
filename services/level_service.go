package services

import (
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/repositories"
)

type LevelService interface {
	ListLevels() ([]*models.Level, error)
	ReadLevel(id uuid.UUID) (*models.Level, error)
	CreateLevel(level *models.Level) error
	UpdateLevel(level *models.Level) error
	DeleteLevel(id uuid.UUID) error
}

type LevelServiceImpl struct {
	Repo repositories.LevelRepository
}

func (s *LevelServiceImpl) ListLevels() ([]*models.Level, error) {
	return s.Repo.ListLevels()
}

func (s *LevelServiceImpl) ReadLevel(id uuid.UUID) (*models.Level, error) {
	return s.Repo.ReadLevel(id)
}

func (s *LevelServiceImpl) CreateLevel(level *models.Level) error {
	level.ID = uuid.Nil
	level.Category.Type = models.Category
	for _, l := range level.LevelNames {
		l.Type = models.LevelName
	}
	return s.Repo.CreateLevel(level)
}

func (s *LevelServiceImpl) UpdateLevel(level *models.Level) error {
	return s.Repo.UpdateLevel(level)
}

func (s *LevelServiceImpl) DeleteLevel(id uuid.UUID) error {
	return s.Repo.DeleteLevel(id)
}

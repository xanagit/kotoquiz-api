package services

import (
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/repositories"
)

type LevelService interface {
	ListLevels() ([]*models.Level, error)
	ReadLevel(id string) (*models.Level, error)
	CreateLevel(label *models.Level) error
	UpdateLevel(label *models.Level) error
	DeleteLevel(id string) error
}

type LevelServiceImpl struct {
	Repo repositories.LevelRepository
}

func (s *LevelServiceImpl) ListLevels() ([]*models.Level, error) {
	return s.Repo.ListLevels()
}

func (s *LevelServiceImpl) ReadLevel(id string) (*models.Level, error) {
	return s.Repo.ReadLevel(id)
}

func (s *LevelServiceImpl) CreateLevel(label *models.Level) error {
	return s.Repo.CreateLevel(label)
}

func (s *LevelServiceImpl) UpdateLevel(label *models.Level) error {
	return s.Repo.UpdateLevel(label)
}

func (s *LevelServiceImpl) DeleteLevel(id string) error {
	return s.Repo.DeleteLevel(id)
}

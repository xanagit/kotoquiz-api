package services

import (
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/repositories"
)

type LabelService interface {
	ListLabels(labelType string) ([]*models.Label, error)
	ReadLabel(id string) (*models.Label, error)
	CreateLabel(label *models.Label) error
	UpdateLabel(label *models.Label) error
	DeleteLabel(id string) error
}

type LabelServiceImpl struct {
	Repo repositories.LabelRepository
}

func (s *LabelServiceImpl) ListLabels(labelType string) ([]*models.Label, error) {
	return s.Repo.ListLabelsByType(labelType)
}

func (s *LabelServiceImpl) ReadLabel(id string) (*models.Label, error) {
	return s.Repo.ReadLabel(id)
}

func (s *LabelServiceImpl) CreateLabel(label *models.Label) error {
	return s.Repo.CreateLabel(label)
}

func (s *LabelServiceImpl) UpdateLabel(label *models.Label) error {
	return s.Repo.UpdateLabel(label)
}

func (s *LabelServiceImpl) DeleteLabel(id string) error {
	return s.Repo.DeleteLabel(id)
}

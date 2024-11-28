package services

import (
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/repositories"
)

type LabelService interface {
	ListLabels(labelType models.LabelType) ([]*models.Label, error)
	ReadLabel(id uuid.UUID) (*models.Label, error)
	CreateLabel(label *models.Label, labelType models.LabelType) error
	UpdateLabel(label *models.Label) error
	DeleteLabel(id uuid.UUID) error
}

type LabelServiceImpl struct {
	Repo repositories.LabelRepository
}

func (s *LabelServiceImpl) ListLabels(labelType models.LabelType) ([]*models.Label, error) {
	return s.Repo.ListLabelsByType(labelType)
}

func (s *LabelServiceImpl) ReadLabel(id uuid.UUID) (*models.Label, error) {
	return s.Repo.ReadLabel(id)
}

func (s *LabelServiceImpl) CreateLabel(label *models.Label, labelType models.LabelType) error {
	label.ID = uuid.Nil
	label.Type = labelType
	return s.Repo.CreateLabel(label)
}

func (s *LabelServiceImpl) UpdateLabel(label *models.Label) error {
	return s.Repo.UpdateLabel(label)
}

func (s *LabelServiceImpl) DeleteLabel(id uuid.UUID) error {
	return s.Repo.DeleteLabel(id)
}

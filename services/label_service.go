package services

import (
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/repositories"
)

type LabelService interface {
	ListLabels(labelType string) ([]*models.Label, error)
	ListLabelsOfCategory(labelType string, categoryId string) ([]*models.Label, error)
	ReadLabel(id string) (*models.Label, error)
	CreateLabel(label *models.Label) error
	CreateCategoryLabel(label *models.Label, categoryId string) error
	UpdateLabel(label *models.Label) error
	DeleteLabel(id string) error
}

type LabelServiceImpl struct {
	Repo repositories.LabelRepository
}

func (s *LabelServiceImpl) ListLabels(labelType string) ([]*models.Label, error) {
	return s.Repo.ListLabelsByType(labelType)
}

func (s *LabelServiceImpl) ListLabelsOfCategory(labelType string, categoryId string) ([]*models.Label, error) {
	return s.Repo.ListLabelsByTypeAndCategory(labelType, categoryId)
}

func (s *LabelServiceImpl) ReadLabel(id string) (*models.Label, error) {
	return s.Repo.ReadLabel(id)
}

func (s *LabelServiceImpl) CreateLabel(label *models.Label) error {
	return s.Repo.CreateLabel(label)
}
func (s *LabelServiceImpl) CreateCategoryLabel(label *models.Label, categoryId string) error {
	return s.Repo.CreateCategoryLabel(label, categoryId)
}

func (s *LabelServiceImpl) UpdateLabel(label *models.Label) error {
	return s.Repo.UpdateLabel(label)
}

func (s *LabelServiceImpl) DeleteLabel(id string) error {
	return s.Repo.DeleteLabel(id)
}

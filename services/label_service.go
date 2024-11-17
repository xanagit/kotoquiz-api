package services

import (
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/repositories"
)

type LabelService interface {
	ListLabels(labelType string) ([]*models.Label, error)
}

type LabelServiceImpl struct {
	Repo repositories.LabelRepository
}

func (s *LabelServiceImpl) ListLabels(labelType string) ([]*models.Label, error) {
	return s.Repo.ListLabelsByType(labelType)
}

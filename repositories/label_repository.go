package repositories

import (
	"github.com/xanagit/kotoquiz-api/models"
	"gorm.io/gorm"
)

type LabelRepository interface {
	ListLabelsByType(labelType string) ([]*models.Label, error)
}

type LabelRepositoryImpl struct {
	DB *gorm.DB
}

func (r *LabelRepositoryImpl) ListLabelsByType(labelType string) ([]*models.Label, error) {
	var labels []*models.Label
	result := r.DB.Where("type = ?", labelType).Find(&labels)
	return labels, result.Error
}

package repositories

import (
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/models"
	"gorm.io/gorm"
)

type LabelRepository interface {
	ListLabelsByType(labelType models.LabelType) ([]*models.Label, error)
	ReadLabel(id uuid.UUID) (*models.Label, error)
	CreateLabel(word *models.Label) error
	UpdateLabel(word *models.Label) error
	DeleteLabel(id uuid.UUID) error
}

type LabelRepositoryImpl struct {
	DB *gorm.DB
}

// Make sure that LabelRepositoryImpl implements LabelRepository
var _ LabelRepository = (*LabelRepositoryImpl)(nil)

func (r *LabelRepositoryImpl) ListLabelsByType(labelType models.LabelType) ([]*models.Label, error) {
	var labels []*models.Label
	result := r.DB.Where("type = ?", labelType).Find(&labels)
	return labels, result.Error
}

func (r *LabelRepositoryImpl) ReadLabel(id uuid.UUID) (*models.Label, error) {
	var label models.Label
	result := r.DB.First(&label, "id = ?", id)
	return &label, result.Error
}

func (r *LabelRepositoryImpl) CreateLabel(label *models.Label) error {
	return r.DB.Create(label).Error
}

func (r *LabelRepositoryImpl) UpdateLabel(label *models.Label) error {
	return r.DB.Save(label).Error
}

func (r *LabelRepositoryImpl) DeleteLabel(id uuid.UUID) error {
	if err := r.DB.Where("id = ?", id).Delete(&models.Label{}).Error; err != nil {
		return err
	}
	return nil
}

package repositories

import (
	"github.com/xanagit/kotoquiz-api/models"
	"gorm.io/gorm"
)

type LabelRepository interface {
	ListLabelsByType(labelType string) ([]*models.Label, error)
	ListLabelsByTypeAndCategory(labelType string, categoryId string) ([]*models.Label, error)
	ReadLabel(id string) (*models.Label, error)
	CreateLabel(word *models.Label) error
	CreateCategoryLabel(label *models.Label, categoryId string) error
	UpdateLabel(word *models.Label) error
	DeleteLabel(id string) error
}

type LabelRepositoryImpl struct {
	DB *gorm.DB
}

func (r *LabelRepositoryImpl) ListLabelsByType(labelType string) ([]*models.Label, error) {
	var labels []*models.Label
	result := r.DB.Where("type = ?", labelType).Find(&labels)
	return labels, result.Error
}

func (r *LabelRepositoryImpl) ListLabelsByTypeAndCategory(labelType string, categoryId string) ([]*models.Label, error) {
	var labels []*models.Label

	result := r.DB.
		Joins("JOIN levels ON labels.id = levels.category_id").
		Where("labels.type = ? AND levels.category_id = ?", labelType, categoryId).
		Find(&labels)

	return labels, result.Error
}

func (r *LabelRepositoryImpl) ReadLabel(id string) (*models.Label, error) {
	var label models.Label
	result := r.DB.First(&label, "id = ?", id)
	return &label, result.Error
}

func (r *LabelRepositoryImpl) CreateLabel(label *models.Label) error {
	return r.DB.Create(label).Error
}

func (r *LabelRepositoryImpl) CreateCategoryLabel(label *models.Label, categoryId string) error {
	// Start a transaction
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Create the label
		if err := tx.Create(label).Error; err != nil {
			return err
		}

		// Find the category
		var category models.Label
		if err := tx.First(&category, "id = ? AND type = ?", categoryId, "CATEGORY").Error; err != nil {
			return err
		}

		// Create the relationship
		if err := tx.Model(&category).Association("LevelNames").Append(label); err != nil {
			return err
		}

		return nil
	})
}

func (r *LabelRepositoryImpl) UpdateLabel(label *models.Label) error {
	return r.DB.Save(label).Error
}

func (r *LabelRepositoryImpl) DeleteLabel(id string) error {
	if err := r.DB.Where("id = ?", id).Delete(&models.Label{}).Error; err != nil {
		return err
	}
	return nil
}

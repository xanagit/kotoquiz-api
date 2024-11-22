package repositories

import (
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/models"
	"gorm.io/gorm"
)

type LevelRepository interface {
	ListLevels() ([]*models.Level, error)
	ReadLevel(id string) (*models.Level, error)
	CreateLevel(word *models.Level) error
	UpdateLevel(word *models.Level) error
	DeleteLevel(id string) error
}

type LevelRepositoryImpl struct {
	DB *gorm.DB
}

func (r *LevelRepositoryImpl) ListLevels() ([]*models.Level, error) {
	var labels []*models.Level
	result := r.DB.Preload("LevelNames").Preload("Category").Find(&labels)
	return labels, result.Error
}

func (r *LevelRepositoryImpl) ReadLevel(id string) (*models.Level, error) {
	var label models.Level
	result := r.DB.First(&label, "id = ?", id)
	return &label, result.Error
}

func (r *LevelRepositoryImpl) CreateLevel(label *models.Level) error {
	return r.DB.Create(label).Error
}

func (r *LevelRepositoryImpl) UpdateLevel(label *models.Level) error {
	return r.DB.Save(label).Error
}

func (r *LevelRepositoryImpl) DeleteLevel(id string) error {
	//if err := r.DB.Where("id = ?", id).Delete(&models.Level{}).Error; err != nil {
	//	return err
	//}
	//return nil
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var level models.Level
		// Charger le level avec ses associations
		if err := tx.First(&level, "id = ?", id).Error; err != nil {
			return err
		}

		// Sauvegarder l'ID de la catégory
		categoryID := level.CategoryID

		// 1. Mettre à null la référence à la catégory
		if err := tx.Model(&level).Update("category_id", nil).Error; err != nil {
			return err
		}

		// 2. Supprimer le level (cela déclenchera BeforeDelete)
		if err := tx.Delete(&level).Error; err != nil {
			return err
		}

		// 3. Maintenant nous pouvons supprimer la catégory
		if categoryID != uuid.Nil {
			if err := tx.Delete(&models.Label{}, "id = ?", categoryID).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

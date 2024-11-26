package repositories

import (
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/models"
	"gorm.io/gorm"
)

type WordRepository interface {
	ListWordsIds(tagIds []string, levelNameIds []string, limit int, offset int) ([]string, error)
	ListWordsByIds(ids []string) ([]*models.Word, error)
	ReadWord(id string) (*models.Word, error)
	CreateWord(word *models.Word) error
	UpdateWord(word *models.Word) error
	DeleteWord(id string) error
}

type WordRepositoryImpl struct {
	DB *gorm.DB
}

func (r *WordRepositoryImpl) ListWordsIds(tagIds []string, levelNameIds []string, limit int, offset int) ([]string, error) {
	var wordIDs []string

	query := r.DB.Table("words w").
		Select("DISTINCT w.id")
	if len(tagIds) > 0 {
		query.
			Joins("JOIN word_tag wt ON wt.word_id = w.id").
			Joins("JOIN labels t ON t.id = wt.label_id").
			Where("t.id IN ?", tagIds)
	}
	if len(levelNameIds) > 0 {
		query.
			Joins("JOIN word_level wl ON wl.word_id = w.id").
			Joins("JOIN level l ON l.id = wl.level_id").
			Joins("JOIN level_values lv ON lv.level_id = l.id").
			Where("lv.label_id IN ?", levelNameIds)
	}
	query.
		Limit(limit).
		Offset(offset)

	err := query.Scan(&wordIDs).Error
	if err != nil {
		return nil, err
	}

	return wordIDs, nil
}

func (r *WordRepositoryImpl) ListWordsByIds(ids []string) ([]*models.Word, error) {
	var words []*models.Word
	// TODO : réécrire la requête pour ne pas utiliser de preload. Utiliser de la pagination ?
	result := r.DB.
		Preload("Tags").
		Preload("Levels.Category").
		Preload("Levels.LevelNames").
		Preload("Translation").
		Where("id IN ?", ids).Find(&words)
	return words, result.Error
}

func (r *WordRepositoryImpl) ReadWord(id string) (*models.Word, error) {
	var word models.Word
	result := r.DB.Preload("Translation").Preload("Tags").Preload("Levels").Preload("Levels.Category").Preload("Levels.LevelNames").First(&word, "id = ?", id)
	return &word, result.Error
}

func (r *WordRepositoryImpl) CreateWord(word *models.Word) error {
	return r.DB.Create(word).Error
}

func (r *WordRepositoryImpl) UpdateWord(word *models.Word) error {
	return r.DB.Save(word).Error
}

func (r *WordRepositoryImpl) DeleteWord(id string) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var word models.Word
		// Charger le word avec ses associations
		if err := tx.First(&word, "id = ?", id).Error; err != nil {
			return err
		}

		// Sauvegarder l'ID de la traduction
		translationID := word.TranslationID

		// 1. Mettre à null la référence à la traduction
		if err := tx.Model(&word).Update("translation_id", nil).Error; err != nil {
			return err
		}

		// 2. Supprimer le word (cela déclenchera BeforeDelete)
		if err := tx.Delete(&word).Error; err != nil {
			return err
		}

		// 3. Maintenant nous pouvons supprimer la traduction
		if translationID != uuid.Nil {
			if err := tx.Delete(&models.Label{}, "id = ?", translationID).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

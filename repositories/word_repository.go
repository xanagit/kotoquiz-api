package repositories

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/models"
	"gorm.io/gorm"
)

type WordRepository interface {
	ListWordsIds(tagIds []string, levelNameIds []string, nb int) ([]string, error)
	ListWordsByIds(ids []uuid.UUID) ([]*models.Word, error)
	ReadWord(id uuid.UUID) (*models.Word, error)
	CreateWord(word *models.Word) error
	UpdateWord(word *models.Word) error
	DeleteWord(id uuid.UUID) error
}

type WordRepositoryImpl struct {
	DB *gorm.DB
}

func (r *WordRepositoryImpl) ListWordsIds(tagIds []string, levelNameIds []string, nb int) ([]string, error) {
	var wordIDs []string

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Define the isolation level for this read
		// REPEATABLE READ prevents modifications while reading
		err := tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ").Error
		if err != nil {
			return err
		}

		query := tx.Table("words w").
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
				Joins("JOIN levels l ON l.id = wl.level_id").
				Joins("JOIN level_values lv ON lv.level_id = l.id").
				Where("lv.label_id IN ?", levelNameIds)
		}
		if nb > 0 {
			query.Limit(nb)
		}

		return query.Scan(&wordIDs).Error
	}, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})

	return wordIDs, err
}

func (r *WordRepositoryImpl) ListWordsByIds(ids []uuid.UUID) ([]*models.Word, error) {
	var words []*models.Word
	result := r.DB.
		Preload("Tags").
		Preload("Levels.Category").
		Preload("Levels.LevelNames").
		Preload("Translation").
		Where("id IN ?", ids).Find(&words)
	return words, result.Error
}

func (r *WordRepositoryImpl) ReadWord(id uuid.UUID) (*models.Word, error) {
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

func (r *WordRepositoryImpl) DeleteWord(id uuid.UUID) error {
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

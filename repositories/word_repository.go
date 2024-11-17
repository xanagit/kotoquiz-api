package repositories

import (
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/models"
	"gorm.io/gorm"
)

type WordRepository interface {
	ListWordsIDsByIds(ids []string, wordIDs *[]uuid.UUID) error
	ListWordsByIds(ids []string) ([]*models.Word, error)
	ReadWord(id string) (*models.Word, error)
	CreateWord(word *models.Word) error
	UpdateWord(word *models.Word) error
	DeleteWord(id string) error
}

type WordRepositoryImpl struct {
	DB *gorm.DB
}

func (r *WordRepositoryImpl) ListWordsIDsByIds(ids []string, wordIDs *[]uuid.UUID) error {
	return r.DB.Model(&models.Word{}).Where("id IN ?", ids).Pluck("id", wordIDs).Error
}

func (r *WordRepositoryImpl) ListWordsByIds(ids []string) ([]*models.Word, error) {
	var words []*models.Word
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
	// .Preload("Tags")
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
	return r.DB.Delete(&models.Word{}, "id = ?", id).Error
}

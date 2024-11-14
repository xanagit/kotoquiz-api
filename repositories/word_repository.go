package repositories

import (
	"github.com/xanagit/kotoquiz-api/models"
	"gorm.io/gorm"
)

type WordRepository interface {
	GetWords() ([]*models.Word, error)
	GetWordByID(id string) (*models.Word, error)
	CreateWord(word *models.Word) error
	UpdateWord(word *models.Word) error
	DeleteWord(id string) error
}

type WordRepositoryImpl struct {
	DB *gorm.DB
}

func (r *WordRepositoryImpl) GetWords() ([]*models.Word, error) {
	var words []*models.Word
	result := r.DB.Find(&words)
	return words, result.Error
}

func (r *WordRepositoryImpl) GetWordByID(id string) (*models.Word, error) {
	var word models.Word
	result := r.DB.First(&word, "id = ?", id)
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

package repositories

import (
	"errors"
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/models"
	"gorm.io/gorm"
)

type WordLearningHistoryRepository interface {
	GetOrInsertHistory(userID uuid.UUID, wordID uuid.UUID) (*models.WordLearningHistory, error)
	CreateHistory(history *models.WordLearningHistory) error
	UpdateHistory(history *models.WordLearningHistory) error
	GetHistoriesByWordIDs(userID uuid.UUID, wordIDs []string) ([]*models.WordLearningHistory, error)
}

type WordLearningHistoryRepositoryImpl struct {
	DB *gorm.DB
}

func (r *WordLearningHistoryRepositoryImpl) GetOrInsertHistory(userID uuid.UUID, wordID uuid.UUID) (*models.WordLearningHistory, error) {
	var history models.WordLearningHistory
	err := r.DB.Where("user_id = ? AND word_id = ?", userID, wordID).First(&history).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Initialize new history if not found
		history = models.WordLearningHistory{
			UserID:         userID,
			WordID:         wordID,
			LearningStatus: models.New,
		}
		err = r.CreateHistory(&history)
	}
	return &history, err
}

func (r *WordLearningHistoryRepositoryImpl) CreateHistory(history *models.WordLearningHistory) error {
	return r.DB.Create(history).Error
}

func (r *WordLearningHistoryRepositoryImpl) UpdateHistory(history *models.WordLearningHistory) error {
	return r.DB.Save(history).Error
}

func (r *WordLearningHistoryRepositoryImpl) GetHistoriesByWordIDs(userID uuid.UUID, wordIDs []string) ([]*models.WordLearningHistory, error) {
	var histories []*models.WordLearningHistory
	err := r.DB.Where("user_id = ? AND word_id IN ?", userID, wordIDs).Find(&histories).Error
	return histories, err
}

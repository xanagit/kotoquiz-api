package repositories

import (
	"errors"
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/models"
	"gorm.io/gorm"
)

type WordLearningHistoryRepository interface {
	GetHistories(userID uuid.UUID, wordIDs []uuid.UUID) (map[uuid.UUID]*models.WordLearningHistory, error)
	InsertHistories(histories []*models.WordLearningHistory) error
	UpdateHistories(histories []*models.WordLearningHistory) error
	GetHistoriesByWordIDs(userID uuid.UUID, wordIDs []string) ([]*models.WordLearningHistory, error)
}

type WordLearningHistoryRepositoryImpl struct {
	DB *gorm.DB
}

func (r *WordLearningHistoryRepositoryImpl) GetHistories(userID uuid.UUID, wordIDs []uuid.UUID) (map[uuid.UUID]*models.WordLearningHistory, error) {
	var histories []*models.WordLearningHistory

	err := r.DB.Set("gorm:query_option", "FOR UPDATE").
		Where("user_id = ? AND word_id IN ?", userID, wordIDs).
		Find(&histories).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create a map of existing histories to make searching easier
	existingHistories := make(map[uuid.UUID]*models.WordLearningHistory)
	for _, h := range histories {
		existingHistories[h.WordID] = h
	}

	// Crée un slice de résultats dans le même ordre que les wordIDs
	//result := make([]*models.WordLearningHistory, len(wordIDs))
	//for i, wordID := range wordIDs {
	//	if history, exists := existingHistories[wordID]; exists {
	//		result[i] = history
	//	} else {
	//		result[i] = nil
	//	}
	//}

	return existingHistories, nil
}

func (r *WordLearningHistoryRepositoryImpl) InsertHistories(histories []*models.WordLearningHistory) error {
	if len(histories) == 0 {
		return nil
	}

	return r.DB.Transaction(func(tx *gorm.DB) error {
		for _, history := range histories {
			// Check existence with lock
			var existing models.WordLearningHistory
			err := tx.Set("gorm:query_option", "FOR UPDATE").
				Where("user_id = ? AND word_id = ?", history.UserID, history.WordID).
				First(&existing).Error

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			// If record not found, insert it
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := tx.Create(history).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *WordLearningHistoryRepositoryImpl) UpdateHistories(histories []*models.WordLearningHistory) error {
	if len(histories) == 0 {
		return nil
	}

	return r.DB.Transaction(func(tx *gorm.DB) error {
		for _, history := range histories {
			// Get a lock on the record before updating
			var existingHistory models.WordLearningHistory
			if err := tx.Set("gorm:query_option", "FOR UPDATE").
				Where("user_id = ? AND word_id = ?", history.UserID, history.WordID).
				First(&existingHistory).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					continue // Skip if history not found
				}
				return err
			}

			// Make update with lock
			if err := tx.Model(&models.WordLearningHistory{}).
				Where("user_id = ? AND word_id = ?", history.UserID, history.WordID).
				Updates(history).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *WordLearningHistoryRepositoryImpl) GetHistoriesByWordIDs(userID uuid.UUID, wordIDs []string) ([]*models.WordLearningHistory, error) {
	var histories []*models.WordLearningHistory
	err := r.DB.Where("user_id = ? AND word_id IN ?", userID, wordIDs).Find(&histories).Error
	return histories, err
}

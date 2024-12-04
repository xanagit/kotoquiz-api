package services

import (
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/dto"
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/repositories"
)

type WordLearningHistoryService interface {
	ProcessQuizResults(userID uuid.UUID, results []dto.WordQuizResult) error
}

type WordLearningHistoryServiceImpl struct {
	Repo repositories.WordLearningHistoryRepository
}

func (s *WordLearningHistoryServiceImpl) ProcessQuizResults(userID uuid.UUID, results []dto.WordQuizResult) error {
	// Build list of word IDs and map of results (word ID -> WordQuizResult)
	wordIDs := make([]uuid.UUID, len(results))
	resultsMap := make(map[uuid.UUID]*dto.WordQuizResult)
	for i, result := range results {
		wordIDs[i] = result.WordID
		resultsMap[result.WordID] = &results[i]
	}

	historiesMap, err := s.Repo.GetHistories(userID, wordIDs)
	if err != nil {
		return err
	}

	var historiesToUpdate []*models.WordLearningHistory
	var historiesToCreate []*models.WordLearningHistory
	for _, result := range results {
		history, exists := historiesMap[result.WordID]
		if !exists {
			history = &models.WordLearningHistory{
				UserID: userID,
				WordID: result.WordID,
			}
			historiesToCreate = append(historiesToCreate, history)
		} else {
			historiesToUpdate = append(historiesToUpdate, history)
		}
		s.updateHistoryBasicInfo(history)
		s.updateHistoryStats(history, resultsMap[history.WordID].Status)
		s.updateLearningStatus(history)
		s.calculateNextReviewDate(history)
	}
	err = s.Repo.UpdateHistories(historiesToUpdate)
	if err != nil {
		return err
	}
	err = s.Repo.InsertHistories(historiesToCreate)
	if err != nil {
		return err
	}

	return nil
}

package services

import (
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/dto"
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/repositories"
	"time"
)

type WordLearningHistoryService interface {
	ProcessQuizResults(userID uuid.UUID, results []dto.WordQuizResult) error
}

type WordLearningHistoryServiceImpl struct {
	Repo repositories.WordLearningHistoryRepository
}

func (s *WordLearningHistoryServiceImpl) ProcessQuizResults(userID uuid.UUID, results []dto.WordQuizResult) error {
	for _, result := range results {
		history, err := s.Repo.GetOrInsertHistory(userID, result.WordID)
		if err != nil {
			return err
		}

		now := time.Now()
		history.LastViewedAt = now
		history.AnswerCount++

		if result.Status == dto.Success {
			history.NbSuccess++
			history.CurrentStreak++
			if history.CurrentStreak > history.BestStreak {
				history.BestStreak = history.CurrentStreak
			}
		} else if result.Status == dto.Error {
			history.NbErrors++
			history.CurrentStreak = 0
		} else if result.Status == dto.Unanswered {
			history.NbUnanswered++
			history.CurrentStreak = 0
		}

		// Update learning status
		s.updateLearningStatus(history)
		// Calculate next review date
		s.calculateNextReviewDate(history)

		if err := s.Repo.UpdateHistory(history); err != nil {
			return err
		}
	}
	return nil
}

func (s *WordLearningHistoryServiceImpl) updateLearningStatus(history *models.WordLearningHistory) {
	totalAnswers := history.NbSuccess + history.NbErrors + history.NbUnanswered
	successRate := float64(history.NbSuccess) / float64(totalAnswers+history.NbUnanswered)

	switch {
	case totalAnswers == 0:
		history.LearningStatus = models.New
	case history.CurrentStreak >= 5 && successRate >= 0.9:
		history.LearningStatus = models.Mastered
	case history.CurrentStreak >= 3 && successRate >= 0.7:
		history.LearningStatus = models.Reviewing
	default:
		history.LearningStatus = models.Learning
	}
}

func (s *WordLearningHistoryServiceImpl) calculateNextReviewDate(history *models.WordLearningHistory) {
	// Base interval according to status
	var baseInterval time.Duration
	switch history.LearningStatus {
	case models.New:
		baseInterval = 4 * time.Hour
	case models.Learning:
		baseInterval = 24 * time.Hour
	case models.Reviewing:
		baseInterval = 72 * time.Hour
	case models.Mastered:
		baseInterval = 168 * time.Hour // 1 week
	}

	// Multiplier factor based on performances
	multiplier := 1.0
	successRate := float64(history.NbSuccess) / float64(history.NbSuccess+history.NbErrors+history.NbUnanswered)

	// Increase interval if good performances
	if history.CurrentStreak > 3 {
		multiplier += float64(history.CurrentStreak) * 0.2 // +20% per streak
	}
	if successRate > 0.8 {
		multiplier += 0.5 // +50% if good success rate
	}

	// Reduce interval if bad performances
	if successRate < 0.6 {
		multiplier -= 0.5 // -50% if bad success rate
	}

	interval := time.Duration(float64(baseInterval) * multiplier)
	history.NextReviewDate = time.Now().Add(interval)
}

package services

import (
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/dto"
	"github.com/xanagit/kotoquiz-api/models"
	"math/rand"
	"sort"
	"time"
)

func (s *WordDtoServiceImpl) prioritizeWords(wordIDs []string, histories []*models.WordLearningHistory, nb int) []string {
	now := time.Now()

	// Create a map for quick access to histories
	historyMap := make(map[string]*models.WordLearningHistory)
	for _, h := range histories {
		historyMap[h.WordID.String()] = h
	}

	// Structure to store word ID with priority
	type wordPriority struct {
		id    string
		score float64
	}

	priorities := make([]wordPriority, 0, len(wordIDs))

	// Compute priority score for each word
	for _, wordID := range wordIDs {
		history, exists := historyMap[wordID]

		var score float64
		if !exists {
			// No history of the word : average score to give a chance to be selected
			score = 50
		} else {
			// Bas score on review date
			timeUntilReview := history.NextReviewDate.Sub(now)

			if timeUntilReview <= 0 {
				// Late review : high priority
				score = 100 + float64(-timeUntilReview.Hours()) // The later, the higher the score is
			} else {
				// Review to come : priority based on date
				score = 100 - float64(timeUntilReview.Hours())
			}

			// Add score according to learning status
			switch history.LearningStatus {
			case models.New:
				score *= 1.2 // Favor new words
			case models.Learning:
				score *= 1.1 // Average high priority for learning in progress words
			case models.Reviewing:
				score *= 0.9 // Average low priority for revision words
			case models.Mastered:
				score *= 0.7 // Low priority for mastered words
			}

			// Bonus for words with low success rate
			if history.AnswerCount > 0 {
				successRate := float64(history.NbSuccess) / float64(history.AnswerCount)
				if successRate < 0.6 {
					score *= 1.3 // Bonus for difficult words
				}
			}
		}

		priorities = append(priorities, wordPriority{id: wordID, score: score})
	}

	// Sort by score descending
	sort.Slice(priorities, func(i, j int) bool {
		return priorities[i].score > priorities[j].score
	})

	// Extract sorted IDs
	result := make([]string, 0, nb)
	for i := 0; i < nb && i < len(priorities); i++ {
		result = append(result, priorities[i].id)
	}

	// Ajouter une légère randomisation pour éviter la prédictibilité
	shuffleTopResults(result)

	return result
}

// Add a slight randomization to results while preserving global priority
func shuffleTopResults(ids []string) {
	// Divide list in groups of priority and shuffle each groups
	groupSize := 10 // Randomization group size
	for i := 0; i < len(ids); i += groupSize {
		end := i + groupSize
		if end > len(ids) {
			end = len(ids)
		}
		group := ids[i:end]
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r.Shuffle(len(group), func(i, j int) {
			group[i], group[j] = group[j], group[i]
		})
	}
}

func shuffleAndLimit(ids []string, limit int) []string {
	// If limit is negative or zero, return an empty list
	if limit <= 0 {
		return []string{}
	}

	// Copy the slice to avoid modifying the original
	shuffled := make([]string, len(ids))
	copy(shuffled, ids)

	// Shuffle the slice using the Fisher-Yates algorithm
	rand.NewSource(time.Now().UnixNano())
	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	//Limit the size of the result
	if len(shuffled) > limit {
		shuffled = shuffled[:limit]
	}

	return shuffled
}

func (s *WordDtoServiceImpl) fetchAndValidateWords(tagIds []string, levelNameIds []string, nb int) (*dto.WordIdsList, error) {
	// Fetch all IDs corresponding to tags and level names
	allWordIDs, err := s.WordRepo.ListWordsIds(tagIds, levelNameIds, -1)
	if err != nil {
		return nil, err
	}

	// If no word matches criteria or nb <= 0, return an empty list
	if len(allWordIDs) == 0 || nb <= 0 {
		return &dto.WordIdsList{Ids: []string{}}, nil
	}

	// Return random selection of words
	return &dto.WordIdsList{Ids: shuffleAndLimit(allWordIDs, nb)}, nil
}

func (s *WordDtoServiceImpl) processWordsWithLearningHistory(userID uuid.UUID, wordIDs []string, nb int) (*dto.WordIdsList, error) {
	// Fetch learning histories
	histories, err := s.LearningHistoryRepo.GetHistoriesByWordIDs(userID, wordIDs)
	if err != nil {
		return nil, err
	}

	// Split words based on history
	withHistory, withoutHistory := s.splitWordsByHistory(wordIDs, histories)

	// Process and combine results
	return s.combineWordLists(withHistory, withoutHistory, histories, nb)
}

func (s *WordDtoServiceImpl) splitWordsByHistory(wordIDs []string, histories []*models.WordLearningHistory) ([]string, []string) {
	historyMap := makeHistoryMap(histories)
	var withHistory, withoutHistory []string

	for _, wordID := range wordIDs {
		if _, exists := historyMap[wordID]; exists {
			withHistory = append(withHistory, wordID)
		} else {
			withoutHistory = append(withoutHistory, wordID)
		}
	}

	return withHistory, withoutHistory
}

func makeHistoryMap(histories []*models.WordLearningHistory) map[string]*models.WordLearningHistory {
	historyMap := make(map[string]*models.WordLearningHistory)
	for _, h := range histories {
		historyMap[h.WordID.String()] = h
	}
	return historyMap
}

func (s *WordDtoServiceImpl) combineWordLists(wordsWithHistory, wordsWithoutHistory []string, histories []*models.WordLearningHistory, nb int) (*dto.WordIdsList, error) {
	// Prioritize words with history
	prioritizedWithHistory := s.prioritizeWords(wordsWithHistory, histories, nb)

	// If we have enough prioritized words, return them
	if len(prioritizedWithHistory) >= nb {
		return &dto.WordIdsList{Ids: prioritizedWithHistory[:nb]}, nil
	}

	// Complete with words without history
	result := s.completeWithWordsWithoutHistory(
		prioritizedWithHistory,
		wordsWithoutHistory,
		nb,
	)

	return &dto.WordIdsList{Ids: result}, nil
}

func (s *WordDtoServiceImpl) completeWithWordsWithoutHistory(prioritizedWords, wordsWithoutHistory []string, nb int) []string {
	remainingNeeded := nb - len(prioritizedWords)

	// Shuffle words without history
	shuffledWords := shuffleAndLimit(wordsWithoutHistory, remainingNeeded)

	// If we still need more words, reuse words with history
	if len(shuffledWords) < remainingNeeded {
		extraWords := s.getRemainingWordsFromHistory(
			prioritizedWords,
			nb,
			remainingNeeded-len(shuffledWords),
		)
		shuffledWords = append(shuffledWords, extraWords...)
	}

	return append(prioritizedWords, shuffledWords...)
}

func (s *WordDtoServiceImpl) getRemainingWordsFromHistory(prioritizedWords []string, nb, stillNeeded int) []string {
	if stillNeeded <= 0 || len(prioritizedWords) <= nb-stillNeeded {
		return []string{}
	}

	extraWords := prioritizedWords[len(prioritizedWords[:nb-stillNeeded]):]
	if len(extraWords) > stillNeeded {
		extraWords = extraWords[:stillNeeded]
	}

	return extraWords
}

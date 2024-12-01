package services

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/dto"
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/repositories"
	"math/rand"
	"sort"
	"time"
)

type WordDtoService interface {
	ListWordsIDs(userID uuid.UUID, tagIds []string, levelNameIds []string, nb int) (*dto.WordIdsList, error)
	ListWordsDtoByIDs(ids []uuid.UUID, lang string) ([]*dto.WordDTO, error)
	ReadWord(id uuid.UUID, lang string) (*dto.WordDTO, error)
}

type WordDtoServiceImpl struct {
	WordRepo            repositories.WordRepository
	LearningHistoryRepo repositories.WordLearningHistoryRepository
}

func (s *WordDtoServiceImpl) ListWordsIDs(userID uuid.UUID, tagIds []string, levelNameIds []string, nb int) (*dto.WordIdsList, error) {
	// 1. Fetch all IDs corresponding to tags and level names
	allWordIDs, err := s.WordRepo.ListWordsIds(tagIds, levelNameIds, -1) // -1 pour récupérer tous les IDs
	if err != nil {
		return nil, err
	}

	// 2. Fetch learning history of those words
	histories, err := s.LearningHistoryRepo.GetHistoriesByWordIDs(userID, allWordIDs)
	if err != nil {
		return nil, err
	}

	// 3. Sort words IDs by priority
	sortedIDs := s.prioritizeWords(allWordIDs, histories, nb)

	return &dto.WordIdsList{Ids: sortedIDs}, nil
}

func (s *WordDtoServiceImpl) ListWordsDtoByIDs(ids []uuid.UUID, lang string) ([]*dto.WordDTO, error) {
	if ids == nil || len(ids) == 0 {
		return nil, fmt.Errorf("no IDs provided")
	}

	// Récupérer les mots correspondant aux IDs
	words, err := s.WordRepo.ListWordsByIds(ids)
	if err != nil {
		return nil, err
	}

	// Mapper les résultats en DTO
	wordDTOs := make([]*dto.WordDTO, len(words))
	for i, word := range words {
		wordDTOs[i] = mapWordToDTO(word, lang)
	}

	return wordDTOs, nil
}

func (s *WordDtoServiceImpl) ReadWord(id uuid.UUID, lang string) (*dto.WordDTO, error) {
	word, err := s.WordRepo.ReadWord(id)
	if err != nil {
		return nil, err
	}

	wordDTO := mapWordToDTO(word, lang)

	return wordDTO, nil
}

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
	groupSize := 3 // Randomization group size
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

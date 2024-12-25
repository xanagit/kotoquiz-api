package services

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/dto"
	"github.com/xanagit/kotoquiz-api/repositories"
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

// Make sure that WordDtoServiceImpl implements WordDtoService
var _ WordDtoService = (*WordDtoServiceImpl)(nil)

func (s *WordDtoServiceImpl) ListWordsIDs(userID uuid.UUID, tagIds []string, levelNameIds []string, nb int) (*dto.WordIdsList, error) {
	// Fetch and validate words
	allWordIDs, err := s.fetchAndValidateWords(tagIds, levelNameIds, nb)
	if err != nil || len(allWordIDs.Ids) == 0 {
		return allWordIDs, err
	}

	// If no user specified, return random words
	if userID == uuid.Nil {
		return allWordIDs, nil
	}

	// Process words with learning history
	return s.processWordsWithLearningHistory(userID, allWordIDs.Ids, nb)
}

func (s *WordDtoServiceImpl) ListWordsDtoByIDs(ids []uuid.UUID, lang string) ([]*dto.WordDTO, error) {
	if ids == nil || len(ids) == 0 {
		return nil, fmt.Errorf("no IDs provided")
	}

	// Fetch words corresponding to IDs
	words, err := s.WordRepo.ListWordsByIds(ids)
	if err != nil {
		return nil, err
	}

	// Map results to DTO
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

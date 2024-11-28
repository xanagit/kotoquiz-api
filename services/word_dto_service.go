package services

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/dto"
	"github.com/xanagit/kotoquiz-api/repositories"
	"math/rand"
	"time"
)

type WordDtoService interface {
	ListWordsIDs(tagIds []string, levelNameIds []string, nb int) (*dto.WordIdsList, error)
	ListWordsDtoByIDs(ids []uuid.UUID, lang string) ([]*dto.WordDTO, error)
	ReadWord(id uuid.UUID, lang string) (*dto.WordDTO, error)
}

type WordDtoServiceImpl struct {
	Repo repositories.WordRepository
}

func (s *WordDtoServiceImpl) ListWordsIDs(tagIds []string, levelNameIds []string, nb int) (*dto.WordIdsList, error) {

	wordIds, err := s.Repo.ListWordsIds(tagIds, levelNameIds, -1) // On récupère tous les IDs
	if err != nil {
		return nil, err
	}
	wordIdsList := dto.WordIdsList{Ids: shuffleAndLimit(wordIds, nb)}
	return &wordIdsList, nil
}

func (s *WordDtoServiceImpl) ListWordsDtoByIDs(ids []uuid.UUID, lang string) ([]*dto.WordDTO, error) {
	if ids == nil || len(ids) == 0 {
		return nil, fmt.Errorf("no IDs provided")
	}

	// Récupérer les mots correspondant aux IDs
	words, err := s.Repo.ListWordsByIds(ids)
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
	word, err := s.Repo.ReadWord(id)
	if err != nil {
		return nil, err
	}

	wordDTO := mapWordToDTO(word, lang)

	return wordDTO, nil
}

func shuffleAndLimit(wordIds []string, nb int) []string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Shuffle the list
	r.Shuffle(len(wordIds), func(i, j int) {
		wordIds[i], wordIds[j] = wordIds[j], wordIds[i]
	})

	// Limit the list to 'nb' items
	if nb > len(wordIds) {
		nb = len(wordIds)
	}
	return wordIds[:nb]
}

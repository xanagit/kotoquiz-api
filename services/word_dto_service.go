package services

import (
	"fmt"
	"github.com/xanagit/kotoquiz-api/dto"
	"github.com/xanagit/kotoquiz-api/repositories"
)

type WordDtoService interface {
	ListWordsIDs(tagIds []string, levelNameIds []string, limit int, offset int) (*dto.WordIdsList, error)
	ListWordsDtoByIDs(ids []string, lang string) ([]*dto.WordDTO, error)
	ReadWord(id string, lang string) (*dto.WordDTO, error)
}

type WordDtoServiceImpl struct {
	Repo repositories.WordRepository
}

func (s *WordDtoServiceImpl) ListWordsIDs(tagIds []string, levelNameIds []string, limit int, offset int) (*dto.WordIdsList, error) {

	wordIds, err := s.Repo.ListWordsIds(tagIds, levelNameIds, limit, offset)
	if err != nil {
		return nil, err
	}
	wordIdsList := dto.WordIdsList{Ids: wordIds}
	return &wordIdsList, nil
}

func (s *WordDtoServiceImpl) ListWordsDtoByIDs(ids []string, lang string) ([]*dto.WordDTO, error) {
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

func (s *WordDtoServiceImpl) ReadWord(id string, lang string) (*dto.WordDTO, error) {
	word, err := s.Repo.ReadWord(id)
	if err != nil {
		return nil, err
	}

	wordDTO := mapWordToDTO(word, lang)

	return wordDTO, nil
}

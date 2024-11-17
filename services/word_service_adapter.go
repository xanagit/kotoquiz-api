package services

import (
	"github.com/xanagit/kotoquiz-api/dto"
	"github.com/xanagit/kotoquiz-api/models"
)

func mapWordToDTO(word *models.Word, lang string) *dto.WordDTO {
	if word == nil {
		return nil
	}

	// Filtrer les Tags en tableau de strings
	var mappedTags []string
	for _, tag := range word.Tags {
		mappedTags = append(mappedTags, extractLabel(tag, lang))
	}

	// Filtrer la Translation en string
	mappedTranslation := extractLabel(&word.Translation, lang)

	// Filtrer les Levels
	var mappedLevels []dto.LevelDTO
	for _, level := range word.Levels {
		// Filtrer la catégorie en string
		mappedCategory := extractLabel(&level.Category, lang)

		// Filtrer les noms des niveaux en tableau de strings
		var mappedLevelNames []string
		for _, levelName := range level.LevelNames {
			mappedLevelNames = append(mappedLevelNames, extractLabel(levelName, lang))
		}

		// Construire un LevelDTO
		mappedLevels = append(mappedLevels, dto.LevelDTO{
			Category:   mappedCategory,
			LevelNames: mappedLevelNames,
		})
	}

	// Construire et retourner un WordDTO
	return &dto.WordDTO{
		ID:          word.ID,
		Kanji:       word.Kanji,
		Onyomi:      word.Onyomi,
		Kunyomi:     word.Kunyomi,
		ImageURL:    word.ImageURL,
		Translation: mappedTranslation,
		Tags:        mappedTags,
		Levels:      mappedLevels,
	}
}

func extractLabel(label *models.Label, lang string) string {
	if label == nil {
		return ""
	}
	switch lang {
	case "en":
		return label.En
	case "fr":
		return label.Fr
	default:
		return ""
	}
}

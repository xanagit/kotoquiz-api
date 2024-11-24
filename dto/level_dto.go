package dto

import (
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/models"
)

type WordDTO struct {
	ID          uuid.UUID       `json:"id"`
	Kanji       string          `json:"kanji"`
	Yomi        string          `json:"yomi"`
	YomiType    models.YomiType `json:"yomiType"`
	ImageURL    string          `json:"image_url"`
	Translation string          `json:"translation"`
	Tags        []string        `json:"tags"`
	Levels      []LevelDTO      `json:"levels"`
}

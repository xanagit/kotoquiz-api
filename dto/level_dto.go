package dto

import "github.com/google/uuid"

type WordDTO struct {
	ID          uuid.UUID  `json:"id"`
	Kanji       string     `json:"kanji"`
	Yomi        string     `json:"yomi"`
	YomiType    string     `json:"yomiType"`
	ImageURL    string     `json:"image_url"`
	Translation string     `json:"translation"`
	Tags        []string   `json:"tags"`
	Levels      []LevelDTO `json:"levels"`
}

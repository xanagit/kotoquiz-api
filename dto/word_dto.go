// Package dto contains Data Transfer Objects used for API communication
// These structures are optimized for sending/receiving data via the API
// and may differ from internal domain models
package dto

import (
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/models"
)

// WordDTO represents a word object optimized for API responses
// It contains only the necessary information for client applications
// and is structured for efficient serialization and deserialization
type WordDTO struct {
	// ID is the unique identifier of the word
	ID          uuid.UUID       `json:"id"`
	Kanji       string          `json:"kanji"`
	Yomi        string          `json:"yomi"`
	YomiType    models.YomiType `json:"yomiType"`
	ImageURL    string          `json:"image_url"`
	Translation string          `json:"translation"`
	Tags        []string        `json:"tags"`
	Levels      []*LevelDTO     `json:"levels"`
}

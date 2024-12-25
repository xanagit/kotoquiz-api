package models

import (
	"github.com/google/uuid"
)

type LabelType string

const (
	Tag         LabelType = "TAG"
	Category    LabelType = "CATEGORY"
	LevelName   LabelType = "LEVEL_NAME"
	Translation LabelType = "TRANSLATION"
)

type Label struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	En   string    `gorm:"size:255" json:"en"`
	Fr   string    `gorm:"size:255" json:"fr"`
	Type LabelType `gorm:"size:100" json:"type"`

	Words []*Word `gorm:"many2many:word_tag;constraint:OnDelete:CASCADE;" json:"-"`
}

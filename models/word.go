package models

import (
	"github.com/google/uuid"
)

type Word struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Kanji         string    `gorm:"size:50" json:"kanji"`
	Onyomi        string    `gorm:"size:50" json:"onyomi"`
	Kunyomi       string    `gorm:"size:50" json:"kunyomi"`
	ImageURL      string    `gorm:"size:255" json:"imageURL"`
	TranslationID uuid.UUID `gorm:"type:uuid" json:"-"`

	Translation Label    `gorm:"foreignKey:TranslationID" json:"translation"`
	Tags        []*Label `gorm:"many2many:word_tag" json:"tags"`
	Levels      []*Level `gorm:"many2many:word_level" json:"levels"`
}

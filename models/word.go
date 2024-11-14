package models

import (
	"github.com/google/uuid"
)

type Word struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Kanji         string    `gorm:"size:50"`
	Onyomi        string    `gorm:"size:50"`
	Kunyomi       string    `gorm:"size:50"`
	ImageURL      string    `gorm:"size:255"`
	TranslationID uuid.UUID `gorm:"type:uuid"`

	Translation Label    `gorm:"foreignKey:TranslationID"`
	Tags        []*Label `gorm:"many2many:word_tag"`
	Levels      []*Label `gorm:"many2many:word_level"`
}

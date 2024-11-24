package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type YomiType string

const (
	Onyomi  YomiType = "ONYOMI"
	Kunyomi YomiType = "KUNYOMI"
)

type Word struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Kanji         string    `gorm:"size:50" json:"kanji"`
	Yomi          string    `gorm:"size:50" json:"yomi"`
	YomiType      YomiType  `gorm:"size:50" json:"yomiType"`
	ImageURL      string    `gorm:"size:255" json:"imageURL"`
	TranslationID uuid.UUID `gorm:"type:uuid" json:"-"`

	Translation Label    `gorm:"foreignKey:TranslationID" json:"translation"`
	Tags        []*Label `gorm:"many2many:word_tag;" json:"tags"`
	Levels      []*Level `gorm:"many2many:word_level;" json:"levels"`
}

// BeforeDelete est un hook GORM qui sera appelé automatiquement avant la suppression
func (w *Word) BeforeDelete(tx *gorm.DB) error {
	// Supprimer les associations many2many
	if err := tx.Model(w).Association("Tags").Clear(); err != nil {
		return err
	}
	if err := tx.Model(w).Association("Levels").Clear(); err != nil {
		return err
	}

	return nil
}

// Package models defines the database entities and their relationships
// It uses GORM for object-relational mapping and includes model hooks and validations
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

// Word represents a vocabulary word in the database
// It contains the core word data and references to translations, tags, and levels
type Word struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Kanji         string    `gorm:"size:50" json:"kanji"`
	Yomi          string    `gorm:"size:50" json:"yomi"`
	YomiType      YomiType  `gorm:"size:50" json:"yomiType"`
	ImageURL      string    `gorm:"size:255" json:"imageURL"`
	TranslationID uuid.UUID `gorm:"type:uuid" json:"-"`

	Translation Label    `gorm:"foreignKey:TranslationID" json:"translation"`
	Tags        []*Label `gorm:"many2many:word_tag;joinForeignKey:WordID;joinReferences:LabelID" json:"tags"`
	Levels      []*Level `gorm:"many2many:word_level;joinForeignKey:WordID;joinReferences:LevelID" json:"levels"`
}

// BeforeDelete is a GORM hook that runs before deleting a word
// It handles cleaning up related records to prevent orphaned relationships
func (w *Word) BeforeDelete(tx *gorm.DB) error {
	// Remove all word-tag associations many2many
	if err := tx.Model(w).Association("Tags").Clear(); err != nil {
		return err
	}
	// Remove all word-level associations
	if err := tx.Model(w).Association("Levels").Clear(); err != nil {
		return err
	}

	return nil
}

package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LevelType string

const (
	BuiltIn     LevelType = "BUILT_IN_LEVEL"
	CustomLevel LevelType = "CUSTOM_LEVEL"
)

type Level struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Type LevelType `gorm:"size:100" json:"type"`
	// TODO : ajouter une fk vers un user s'il s'agit de la liste d'un user
	CategoryID uuid.UUID `gorm:"type:uuid" json:"-"`

	Category   Label    `gorm:"foreignKey:CategoryID" json:"category"`
	LevelNames []*Label `gorm:"many2many:level_values;constraint:OnDelete:CASCADE;" json:"levelNames"`
	Words      []*Word  `gorm:"many2many:word_level" json:"-"`
}

// BeforeDelete est un hook GORM qui sera appelé automatiquement avant la suppression
func (l *Level) BeforeDelete(tx *gorm.DB) error {
	// Supprimer les associations many2many
	if err := tx.Model(l).Association("LevelNames").Clear(); err != nil {
		return err
	}
	if err := tx.Model(l).Association("Words").Clear(); err != nil {
		return err
	}

	return nil
}

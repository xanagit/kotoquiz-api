package models

import (
	"github.com/google/uuid"
)

type Label struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	En string    `gorm:"size:255"`
	Fr string    `gorm:"size:255"`

	Words      []*Word `gorm:"many2many:word_tag"`
	WordLevels []*Word `gorm:"many2many:word_level"`
}

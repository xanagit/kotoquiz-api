package models

import (
	"github.com/google/uuid"
)

type Level struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CategoryID uuid.UUID `gorm:"type:uuid" json:"-"`

	Category   Label    `gorm:"foreignKey:CategoryID" json:"category"`
	LevelNames []*Label `gorm:"many2many:level_values" json:"levelNames"`
	Words      []*Word  `gorm:"many2many:word_level" json:"-"`
}

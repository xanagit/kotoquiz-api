package models

import "github.com/google/uuid"

type WordTag struct {
	WordID  uuid.UUID `gorm:"type:uuid;primaryKey;index:idx_word_tag,priority:1"`
	LabelID uuid.UUID `gorm:"type:uuid;primaryKey;index:idx_word_tag,priority:2"`
}

type WordLevel struct {
	WordID  uuid.UUID `gorm:"type:uuid;primaryKey;index:idx_word_level,priority:1"`
	LevelID uuid.UUID `gorm:"type:uuid;primaryKey;index:idx_word_level,priority:2"`
}

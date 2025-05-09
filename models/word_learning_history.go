package models

import (
	"github.com/google/uuid"
	"time"
)

type WLStatus string

const (
	New       WLStatus = "NEW"
	Learning  WLStatus = "LEARNING"
	Reviewing WLStatus = "REVIEWING"
	Mastered  WLStatus = "MASTERED"
)

type WordLearningHistory struct {
	UserID string    `gorm:"type:varchar(255);primaryKey;index:idx_user_word,priority:1" json:"userId"` // Keycloak user ID
	WordID uuid.UUID `gorm:"type:uuid;primaryKey;index:idx_user_word,priority:2" json:"wordId"`

	// Learning timing
	LastViewedAt   time.Time `json:"lastViewedAt"`
	NextReviewDate time.Time `json:"nextReviewDate"`

	// Learning statistics
	AnswerCount   int `gorm:"default:0" json:"viewCount"`
	NbSuccess     int `gorm:"default:0" json:"nbSuccess"`
	NbErrors      int `gorm:"default:0" json:"nbErrors"`
	NbUnanswered  int `gorm:"default:0" json:"nbUnanswered"`
	CurrentStreak int `gorm:"default:0" json:"currentStreak"`
	BestStreak    int `gorm:"default:0" json:"bestStreak"`

	// Learning Status
	LearningStatus WLStatus `gorm:"default:'NEW'" json:"learningStatus"`

	// Relations
	Word Word `gorm:"foreignKey:WordID" json:"-"`
}

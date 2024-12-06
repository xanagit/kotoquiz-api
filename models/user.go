package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email     string    `gorm:"size:255;unique;not null" json:"email"`
	Username  string    `gorm:"size:50;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"` // Don't return password in JSON
	Role      UserRole  `gorm:"size:50;not null" json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Relations
	LearningHistories []*WordLearningHistory `json:"learningHistories,omitempty"`
}

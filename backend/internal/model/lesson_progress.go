package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LessonProgress struct {
	ID           string         `gorm:"primaryKey" json:"id"`
	UserID       string         `gorm:"uniqueIndex:idx_user_lesson;not null" json:"user_id"`
	LessonID     string         `gorm:"uniqueIndex:idx_user_lesson;not null" json:"lesson_id"`
	IsCompleted  bool           `gorm:"default:false" json:"is_completed"`
	LastPosition int            `gorm:"default:0" json:"last_position"` // video resume in seconds
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Associations
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Lesson Lesson `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
}

func (lp *LessonProgress) BeforeCreate(tx *gorm.DB) error {
	if lp.ID == "" {
		lp.ID = uuid.New().String()
	}
	return nil
}

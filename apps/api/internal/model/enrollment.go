package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Enrollment struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	UserID      string         `gorm:"uniqueIndex:idx_user_course;not null" json:"user_id"`
	CourseID    string         `gorm:"uniqueIndex:idx_user_course;not null" json:"course_id"`
	EnrolledAt  time.Time      `gorm:"default:now()" json:"enrolled_at"`
	CompletedAt *time.Time     `json:"completed_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Associations
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Course Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

func (e *Enrollment) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return nil
}

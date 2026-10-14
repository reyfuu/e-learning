package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	UserID    string         `gorm:"uniqueIndex:idx_user_course_review;not null" json:"user_id"`
	CourseID  string         `gorm:"uniqueIndex:idx_user_course_review;not null" json:"course_id"`
	Rating    int            `gorm:"not null" json:"rating"` // check rating between 1 and 5
	Comment   *string        `json:"comment"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Associations
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Course Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Section struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	CourseID  string         `gorm:"index;not null" json:"course_id"`
	Title     string         `gorm:"not null" json:"title"`
	Position  int            `gorm:"not null" json:"position"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Associations
	Lessons []Lesson `gorm:"foreignKey:SectionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"lessons,omitempty"`
}

func (s *Section) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

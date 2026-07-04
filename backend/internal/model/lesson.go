package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Lesson struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	SectionID   string         `gorm:"index;not null" json:"section_id"`
	Title       string         `gorm:"not null" json:"title"`
	ContentType string         `gorm:"type:varchar(20);not null" json:"content_type"` // video | text | quiz | file
	VideoURL    *string        `json:"video_url"`
	Content     *string        `json:"content"`
	Duration    int            `json:"duration"` // in seconds
	Position    int            `gorm:"not null" json:"position"`
	IsPreview   bool           `gorm:"default:false" json:"is_preview"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (l *Lesson) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	return nil
}

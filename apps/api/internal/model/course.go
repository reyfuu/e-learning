package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Course struct {
	ID            string         `gorm:"primaryKey" json:"id"`
	UserID        string         `gorm:"column:instructor_id;index;not null" json:"user_id"` // Map struct's UserID to instructor_id
	Title         string         `gorm:"not null" json:"title"`
	Link          string         `gorm:"column:link" json:"link"` // Keep for compatibility
	Icon          string         `gorm:"column:icon;default:'book'" json:"icon"` // Keep for compatibility
	Slug          string         `gorm:"uniqueIndex;not null" json:"slug"`
	Description   string         `json:"description"`
	ThumbnailURL  string         `gorm:"column:thumbnail_url" json:"thumbnail_url"`
	Price         float64        `gorm:"type:numeric(12,2);default:0" json:"price"`
	Level         string         `json:"level"` // beginner | intermediate | advanced
	Status        string         `gorm:"default:'draft'" json:"status"` // draft | pending | published | archived | learning | completed | cancelled
	CategoryID    *string        `gorm:"column:category_id;index" json:"category_id"`
	Language      string         `gorm:"default:'id'" json:"language"`
	TotalDuration int            `gorm:"column:total_duration;default:0" json:"total_duration"` // in seconds
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Associations
	Instructor  User         `gorm:"foreignKey:UserID" json:"instructor,omitempty"`
	Category    *Category    `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Sections    []Section    `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"sections,omitempty"`
	Enrollments []Enrollment `gorm:"foreignKey:CourseID" json:"enrollments,omitempty"`
	Reviews     []Review     `gorm:"foreignKey:CourseID" json:"reviews,omitempty"`
}

func (c *Course) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

// Request and Response Types

type CreateCourseRequest struct {
	Title string `json:"title" validate:"required"`
	Link  string `json:"link"`
	Icon  string `json:"icon"`
}

type CreateModuleRequest struct {
	Title    string `json:"title" validate:"required"`
	Duration string `json:"duration"`
}

type UpdateCourseStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

type CourseResponse struct {
	ID               string           `json:"id"`
	Title            string           `json:"title"`
	Link             string           `json:"link"`
	Icon             string           `json:"icon"`
	Status           string           `json:"status"`
	Progress         int              `json:"progress"`
	TotalModules     int              `json:"total_modules"`
	CompletedModules int              `json:"completed_modules"`
	Modules          []ModuleResponse `json:"modules"`
}

type ModuleResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	IsCompleted bool   `json:"is_completed"`
	Duration    string `json:"duration"`
}

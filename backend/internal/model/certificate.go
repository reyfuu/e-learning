package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Certificate struct {
	ID         string         `gorm:"primaryKey" json:"id"`
	UserID     string         `gorm:"uniqueIndex:idx_user_course_cert;not null" json:"user_id"`
	CourseID   string         `gorm:"uniqueIndex:idx_user_course_cert;not null" json:"course_id"`
	IssuedAt   time.Time      `gorm:"default:now()" json:"issued_at"`
	CertURL    *string        `json:"cert_url"`
	VerifyCode string         `gorm:"uniqueIndex;not null" json:"verify_code"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Associations
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Course Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

func (c *Certificate) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID            string         `gorm:"primaryKey" json:"id"`
	UserID        string         `gorm:"index;not null" json:"user_id"`
	CourseID      string         `gorm:"index;not null" json:"course_id"`
	Amount        float64        `gorm:"type:numeric(12,2);not null" json:"amount"`
	PlatformFee   float64        `gorm:"type:numeric(12,2);not null" json:"platform_fee"`
	Status        string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending | success | failed | refunded
	PaymentMethod *string        `gorm:"type:varchar(50)" json:"payment_method"`
	GatewayRef    *string        `gorm:"type:varchar(255)" json:"gateway_ref"`
	PaidAt        *time.Time     `json:"paid_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Associations
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Course Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

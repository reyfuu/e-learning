package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         string         `gorm:"primaryKey" json:"id"`
	Email      string         `gorm:"uniqueIndex;not null" json:"email"`
	Password   string         `gorm:"not null" json:"-"`
	Name       string         `gorm:"column:full_name;not null" json:"name"`
	AvatarURL  *string        `gorm:"column:avatar_url" json:"avatar_url"`
	Role       string         `gorm:"column:role;default:'learner'" json:"role"` // learner | instructor | admin
	IsVerified bool           `gorm:"column:is_verified;default:false" json:"is_verified"`
	IsActive   bool           `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
}

type UpdateUserRequest struct {
	Name      string  `json:"name"`
	AvatarURL *string `json:"avatar_url"`
}

type UserResponse struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	AvatarURL  *string   `json:"avatar_url"`
	Role       string    `json:"role"`
	IsVerified bool      `json:"is_verified"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:         u.ID,
		Email:      u.Email,
		Name:       u.Name,
		AvatarURL:  u.AvatarURL,
		Role:       u.Role,
		IsVerified: u.IsVerified,
		IsActive:   u.IsActive,
		CreatedAt:  u.CreatedAt,
	}
}

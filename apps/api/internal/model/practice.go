package model

import "time"

type PracticeProgress struct {
	UserID        string    `gorm:"primaryKey;size:36" json:"-"`
	ExerciseID    string    `gorm:"primaryKey;size:80" json:"exercise_id"`
	HintsUnlocked int       `gorm:"not null;default:0" json:"hints_unlocked"`
	UpdatedAt     time.Time `json:"updated_at"`
}

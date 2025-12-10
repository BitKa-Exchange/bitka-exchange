package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key"`
	Email        string    `gorm:"uniqueIndex;not null"`
	Username     string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserRegisterEvent struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
}
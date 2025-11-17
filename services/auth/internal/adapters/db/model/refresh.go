package model

import (
	"time"
)

// GORM model for the auth service refresh token.
// Put in service module (service-specific).
type RefreshToken struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	TokenHash string    `gorm:"size:128;uniqueIndex;not null"`
	UserID    uint64    `gorm:"not null;index"`
	ExpiresAt time.Time `gorm:"not null;index"`
	Revoked   bool      `gorm:"not null;default:false;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

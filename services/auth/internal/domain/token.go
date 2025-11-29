package domain

import (
	"time"

	"github.com/google/uuid"
)

// TODO: Move these to pkg/token

// RefreshToken represents the record stored in DB for revocation
type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `gorm:"index"`
	TokenJTI  string    `gorm:"uniqueIndex"`
	ExpiresAt time.Time
	IsRevoked bool `gorm:"default:false"`
}

// TokenPair is a Value Object returned by Usecase
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

package domain

import (
	"time"

	"github.com/google/uuid"
)

// Profile is specific to Account Service.
// It shares the same ID as Auth User, but lives in a different DB.
type Profile struct {
	UserID    uuid.UUID `gorm:"type:uuid;primary_key"`
	FullName  string
	AvatarURL string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AccountRepository interface {
	GetProfile(userID uuid.UUID) (*Profile, error)
	UpsertProfile(profile *Profile) error
}

type AccountUsecase interface {
	GetMyProfile(userID string) (*Profile, error)
	UpdateMyProfile(userID string, fullName, avatar string) error
}

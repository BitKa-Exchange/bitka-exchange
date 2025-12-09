package domain

import (
	"time"

	"github.com/google/uuid"
)

// Profile is specific to Account Service.
// It shares the same ID as Auth User, but lives in a different DB.
type Profile struct {
	UserID    uuid.UUID `gorm:"type:uuid;primary_key"`
	Email     string
	FirstName string
	LastName  string
	Username  string
	Age       int
	AvatarURL string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AccountRepository interface {
	GetProfile(userID uuid.UUID) (*Profile, error)
	CreateProfile(userID uuid.UUID, email string , username string) error
	UpsertProfile(profile *Profile) error
}

type AccountUsecase interface {
	GetMyProfile(userID uuid.UUID) (*Profile, error)
	UpdateMyProfile(userID uuid.UUID, fullName, avatar string) error
	CreateUserProfile(userID uuid.UUID, email, username string) error
}

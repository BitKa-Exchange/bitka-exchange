package repository

import (
	"bitka/services/account/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type accountRepo struct {
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) domain.AccountRepository {
	return &accountRepo{db: db}
}

func (r *accountRepo) CreateProfile(userID uuid.UUID, email, username string) error {
	p := domain.Profile{
		UserID:    userID,
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return r.db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			UpdateAll: true,
		}).
		Create(&p).Error
}

func (r *accountRepo) GetProfile(userID uuid.UUID) (*domain.Profile, error) {
	var profile domain.Profile
	err := r.db.First(&profile, "user_id = ?", userID).Error
	return &profile, err
}

func (r *accountRepo) UpsertProfile(profile *domain.Profile) error {
	// Upsert: Create or Update if exists
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"full_name", "avatar_url", "updated_at"}),
	}).Create(profile).Error
}

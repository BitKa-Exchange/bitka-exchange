package db

import (
	"errors"
	"time"

	// import your libs/db module
	"gorm.io/gorm"

	"bitka/auth-service/internal/adapters/db/model"
)

// GormRefreshRepo implements model.RefreshTokenRepository using GORM.
type GormRefreshRepo struct {
	gdb *gorm.DB
}

// NewGormRefreshRepo returns the repo. gdb should already have search_path set to service schema.
func NewGormRefreshRepo(gdb *gorm.DB) *GormRefreshRepo {
	return &GormRefreshRepo{gdb: gdb}
}

func (r *GormRefreshRepo) Save(hashed string, userID uint64, expiresAt time.Time) error {
	if hashed == "" {
		return errors.New("hashed required")
	}
	rec := &model.RefreshToken{
		TokenHash: hashed,
		UserID:    userID,
		ExpiresAt: expiresAt.UTC(),
		Revoked:   false,
	}
	if err := r.gdb.Create(rec).Error; err != nil {
		// handle unique violation if desired; here we bubble up
		return err
	}
	return nil
}

func (r *GormRefreshRepo) Find(hashed string) (uint64, time.Time, bool, error) {
	if hashed == "" {
		return 0, time.Time{}, false, errors.New("hashed required")
	}
	var rec model.RefreshToken
	err := r.gdb.Where("token_hash = ? AND revoked = FALSE", hashed).Take(&rec).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, time.Time{}, false, nil
	}
	if err != nil {
		return 0, time.Time{}, false, err
	}
	return rec.UserID, rec.ExpiresAt, true, nil
}

func (r *GormRefreshRepo) Delete(hashed string) error {
	if hashed == "" {
		return errors.New("hashed required")
	}
	if err := r.gdb.Where("token_hash = ?", hashed).Delete(&model.RefreshToken{}).Error; err != nil {
		return err
	}
	return nil
}

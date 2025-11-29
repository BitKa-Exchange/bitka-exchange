package usecase

import (
	"errors"
	"time"

	"bitka/services/account/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type accountUC struct {
	repo domain.AccountRepository
}

func NewAccountUsecase(repo domain.AccountRepository) domain.AccountUsecase {
	return &accountUC{repo: repo}
}

func (u *accountUC) GetMyProfile(userIDStr string) (*domain.Profile, error) {
	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	profile, err := u.repo.GetProfile(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return a default empty profile if not found, or error based on requirements
			return &domain.Profile{UserID: uid}, nil
		}
		return nil, err
	}
	return profile, nil
}

func (u *accountUC) UpdateMyProfile(userIDStr, name, avatar string) error {
	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		return errors.New("invalid user id")
	}

	return u.repo.UpsertProfile(&domain.Profile{
		UserID:    uid,
		FullName:  name,
		AvatarURL: avatar,
		UpdatedAt: time.Now(),
	})
}

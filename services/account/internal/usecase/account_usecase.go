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

func (u *accountUC) CreateUserProfile (id uuid.UUID ,email,username string) error {
	return u.repo.CreateProfile(id, email, username)
}

func (u *accountUC) GetMyProfile(id uuid.UUID) (*domain.Profile, error) {
	profile, err := u.repo.GetProfile(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return a default empty profile if not found, or error based on requirements
			return &domain.Profile{UserID: id}, nil
		}
		return nil, err
	}
	return profile, nil
}

func (u *accountUC) UpdateMyProfile(id uuid.UUID, fullName, avatar string) error {
	profile := domain.Profile{
        UserID: id,
        FullName: fullName,     
        AvatarURL: avatar,
        UpdatedAt: time.Now(),
    }

    return u.repo.UpsertProfile(&profile)
}

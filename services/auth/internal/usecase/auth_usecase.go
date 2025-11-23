package usecase

import (
	"errors"
	"time"

	"bitka/services/auth/internal/domain"

	"github.com/google/uuid"
)

type authUsecase struct {
	repo     domain.AuthRepository
	tokenGen domain.TokenGenerator
}

func NewAuthUsecase(repo domain.AuthRepository, tg domain.TokenGenerator) domain.AuthUsecase {
	return &authUsecase{repo: repo, tokenGen: tg}
}

func (u *authUsecase) Login(email, password string) (*domain.TokenPair, error) {
	user, err := u.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// In real app: bcrypt.CompareHashAndPassword
	if user.PasswordHash != password {
		return nil, errors.New("invalid credentials")
	}

	// 1. Access Token (15 mins)
	access, err := u.tokenGen.Generate(user.ID.String(), 15*time.Minute, "api:access", "")
	if err != nil {
		return nil, err
	}

	// 2. Refresh Token (7 days)
	refreshJTI := uuid.New().String()
	refresh, err := u.tokenGen.Generate(user.ID.String(), 7*24*time.Hour, "api:refresh", refreshJTI)
	if err != nil {
		return nil, err
	}

	// 3. Persist Refresh Token
	err = u.repo.SaveRefreshToken(&domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenJTI:  refreshJTI,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (u *authUsecase) Register(email, password string) error {
	// In real app: Hash password
	user := &domain.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: password, // TODO: Hash this!
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return u.repo.CreateUser(user)
}

func (u *authUsecase) GetJWKS() ([]byte, error) {
	return u.tokenGen.GetJWKS()
}

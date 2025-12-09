package usecase

import (
	"errors"
	"time"
	"log"

	"bitka/services/auth/internal/domain"
	"bitka/services/auth/internal/kafka"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	repo          domain.AuthRepository
	tokenGen      domain.TokenGenerator
	kafkaProducer *kafka.Producer
}

func NewAuthUsecase(
	repo domain.AuthRepository,
	tg domain.TokenGenerator,
	kp *kafka.Producer,
) domain.AuthUsecase {
	return &authUsecase{
		repo:          repo,
		tokenGen:      tg,
		kafkaProducer: kp,
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (u *authUsecase) Login(identifier, password string) (*domain.TokenPair, error) {
	user, err := u.repo.FindByEmailOrUser(identifier)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !checkPassword(password, user.PasswordHash) {
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

func (u *authUsecase) Register(email, username, password string) error {
	hash_password, err := hashPassword(password)
	if err != nil {
		return err
	}
	user := &domain.User{
		ID:           uuid.New(),
		Email:        email,
		Username:     username,
		PasswordHash: hash_password,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := u.repo.CreateUser(user); err != nil {
		return err
	}
	event := kafka.UserRegisterEvent{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
	}
	err = u.kafkaProducer.PublishUserRegister(event)
	if err != nil {
		log.Println("Failed to publish Kafka event:", err)
		return err
	}
	return u.kafkaProducer.PublishUserRegister(event)
}

func (u *authUsecase) GetJWKS() ([]byte, error) {
	return u.tokenGen.GetJWKS()
}

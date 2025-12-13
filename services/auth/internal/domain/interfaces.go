package domain

import "time"

// TODO: Decide if these interface should be here or not.

// AuthRepository defines data access methods
type AuthRepository interface {
	CreateUser(user *User) error
	FindByEmailOrUser(identifier string) (*User, error)
	SaveRefreshToken(token *RefreshToken) error
}

// AuthUsecase defines business logic methods
type AuthUsecase interface {
	Login(email, password string) (*TokenPair, error)
	Register(email, username, password string) error
	GetJWKS() ([]byte, error)
}

// TokenGenerator defines the behavior we need from pkg/token
// This allows us to mock the complex JWX library in tests
type TokenGenerator interface {
	Generate(userID string, duration time.Duration, audience string, jti string) (string, error)
	GetJWKS() ([]byte, error)
}

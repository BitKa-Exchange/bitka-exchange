package db

import "time"

// RefreshTokenRepository is the interface for storing refresh tokens.
// Implement with DB in production. For demo we provide an in-memory impl.
type RefreshTokenRepository interface {
	Save(hashed string, userID uint64, expiresAt time.Time) error
	Find(hashed string) (userID uint64, expiresAt time.Time, ok bool, err error)
	Delete(hashed string) error
}

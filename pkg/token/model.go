package token

import (
	"time"
)

// RSAKey represents a rotated key pair stored in the database.
type RSAKey struct {
	KID        string `gorm:"primaryKey"` // The unique Key ID
	Algorithm  string `gorm:"size:10"`    // e.g., RS256
	PublicPEM  []byte `gorm:"type:text"`  // PEM encoded Public Key
	PrivatePEM []byte `gorm:"type:text"`  // PEM encoded Private Key (Keep this DB secure!)
	CreatedAt  time.Time
	ExpiresAt  time.Time // Optional: When this key should stop being used for signing
}

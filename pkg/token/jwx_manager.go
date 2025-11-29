package token

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jws"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"gorm.io/gorm"
)

// Manager handles key rotation and signing with DB persistence.
type Manager struct {
	mu           sync.RWMutex
	db           *gorm.DB
	privateKey   *rsa.PrivateKey
	currentKeyID string
	publicKey    jwk.Key
}

// NewManager initializes the manager and syncs with the Database.
// TODO: Use KMS for production use.
// TODO: Look like the key ID isn't used correctly for key set.
func NewManager(db *gorm.DB) (*Manager, error) {
	// 1. Ensure the keys table exists
	if err := db.AutoMigrate(&RSAKey{}); err != nil {
		return nil, err
	}

	m := &Manager{db: db}

	// 2. Try to load the most recent key from DB
	var latestKey RSAKey
	err := db.Order("created_at desc").First(&latestKey).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("No keys found in DB. Generating initial key pair...")
			if err := m.Rotate(); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		log.Printf("Loaded existing key pair (KID: %s) from DB", latestKey.KID)
		if err := m.loadKeyFromStruct(latestKey); err != nil {
			return nil, err
		}
	}

	return m, nil
}

// Rotate generates a new key, saves it to DB, and updates memory.
func (m *Manager) Rotate() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 1. Generate RSA Key
	rawPriv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// 2. Prepare PEM data for DB storage
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(rawPriv),
	})

	rawPub := &rawPriv.PublicKey
	pubBytes, err := x509.MarshalPKIXPublicKey(rawPub)
	if err != nil {
		return err
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})

	// 3. Create Metadata
	kid := uuid.New().String()

	dbKey := RSAKey{
		KID:        kid,
		Algorithm:  "RS256",
		PrivatePEM: privPEM,
		PublicPEM:  pubPEM,
		CreatedAt:  time.Now(),
		// Optional: Set expiry for rotation policy (e.g., 30 days)
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	// 4. Save to Database
	if err := m.db.Create(&dbKey).Error; err != nil {
		return err
	}

	// 5. Update Memory
	return m.loadKeyFromStruct(dbKey)
}

// loadKeyFromStruct parses the DB model into usable in-memory keys
func (m *Manager) loadKeyFromStruct(k RSAKey) error {
	// Parse Private Key
	privBlock, _ := pem.Decode(k.PrivatePEM)
	if privBlock == nil {
		return errors.New("failed to decode private key pem")
	}
	rawPriv, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		return err
	}

	// Create JWK Public Key
	pubKey, err := jwk.PublicKeyOf(rawPriv)
	if err != nil {
		return err
	}

	// These are critical for the Validator to accept the key
	if err := pubKey.Set(jwk.KeyIDKey, k.KID); err != nil {
		return err
	}
	if err := pubKey.Set(jwk.AlgorithmKey, jwa.RS256()); err != nil {
		return err
	}
	// "use": "sig" tells validators this key is for signatures
	if err := pubKey.Set(jwk.KeyUsageKey, jwk.ForSignature); err != nil {
		return err
	}

	m.privateKey = rawPriv
	m.currentKeyID = k.KID
	m.publicKey = pubKey

	return nil
}

// Generate creates a signed JWT string.
func (m *Manager) Generate(userID string, duration time.Duration, audience string, jti string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	builder := jwt.NewBuilder().
		Issuer("bitka-auth").
		Subject(userID).
		Audience([]string{audience}).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(duration))

	if jti != "" {
		builder.JwtID(jti)
	}

	token, err := builder.Build()
	if err != nil {
		return "", err
	}

	headers := jws.NewHeaders()
	headers.Set(jws.KeyIDKey, m.currentKeyID)
	headers.Set(jws.TypeKey, "JWT")

	// Pass headers inside jwt.WithKey()
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256(), m.privateKey, jws.WithProtectedHeaders(headers)))

	if err != nil {
		return "", err
	}

	return string(signed), nil
}

// GetJWKS returns the JSON Web Key Set.
// Note: In a real system, you might want to query ALL valid keys from the DB
// to populate the JWKS, so clients can verify tokens signed by older (but valid) keys.
// For now, we return the current active key.
func (m *Manager) GetJWKS() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	set := jwk.NewSet()
	set.AddKey(m.publicKey)

	return json.Marshal(set)
}

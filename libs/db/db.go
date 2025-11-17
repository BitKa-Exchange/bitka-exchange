package db

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// RefreshToken is the persistent model for refresh tokens.
// - TokenHash should be unique (we store hash, not plain token).
// - We store expires_at and a revoked flag for quick checks.
// - CreatedAt/UpdatedAt are filled by GORM.
type RefreshToken struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	TokenHash string    `gorm:"size:128;uniqueIndex;not null"` // base64(SHA256) ~ 43 chars
	UserID    uint64    `gorm:"not null;index"`                // user foreign id (not FK enforced here)
	ExpiresAt time.Time `gorm:"not null;index"`
	Revoked   bool      `gorm:"not null;default:false;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewGormDB opens a GORM DB with sensible defaults for Postgres/pgx.
// DSN example: "postgres://user:pass@localhost:5432/bitka?sslmode=disable"
func NewGormDB(dsn string, debug bool) (*gorm.DB, error) {
	if dsn == "" {
		return nil, errors.New("dsn required")
	}

	// Configure GORM logger level (Silent by default)
	var gormLogger logger.Interface
	if debug {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: false, // pgx: false uses extended protocol
	})

	gormCfg := &gorm.Config{
		Logger: gormLogger,
		// Add other GORM config here if needed
	}

	db, err := gorm.Open(dialector, gormCfg)
	if err != nil {
		return nil, fmt.Errorf("gorm open: %w", err)
	}

	// configure underlying sql.DB connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("db.DB(): %w", err)
	}
	// sensible defaults — tune per workload
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	// ping to verify connectivity
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	return db, nil
}

// AutoMigrate runs GORM automigration for RefreshToken model.
// Call this during bootstrap/migrations step.
func AutoMigrate(db *gorm.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}
	return db.AutoMigrate(&RefreshToken{})
}

// GormRefreshRepo implements a simple repository for refresh tokens using GORM.
type GormRefreshRepo struct {
	db *gorm.DB
}

// NewGormRefreshRepo constructs the repo. db must be non-nil.
func NewGormRefreshRepo(db *gorm.DB) *GormRefreshRepo {
	return &GormRefreshRepo{db: db}
}

// Save inserts a new refresh token record. If a duplicate token hash exists,
// it returns nil (idempotent behavior) — practical because hashed token
// collision should not happen but Save may be retried.
func (r *GormRefreshRepo) Save(hashed string, userID uint64, expiresAt time.Time) error {
	if hashed == "" {
		return errors.New("hashed token required")
	}
	rec := &RefreshToken{
		TokenHash: hashed,
		UserID:    userID,
		ExpiresAt: expiresAt.UTC(),
		Revoked:   false,
	}
	// Use Create; handle unique constraint error gracefully
	result := r.db.Create(rec)
	if result.Error != nil {
		// if unique violation, treat as success (idempotent)
		if isUniqueViolation(result.Error) {
			return nil
		}
		return result.Error
	}
	return nil
}

// Find looks up a refresh token by its hashed value. Returns userID, expiresAt, found, error.
func (r *GormRefreshRepo) Find(hashed string) (uint64, time.Time, bool, error) {
	if hashed == "" {
		return 0, time.Time{}, false, errors.New("hashed token required")
	}
	var rec RefreshToken
	res := r.db.Where("token_hash = ?", hashed).Take(&rec)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return 0, time.Time{}, false, nil
	}
	if res.Error != nil {
		return 0, time.Time{}, false, res.Error
	}
	return rec.UserID, rec.ExpiresAt, true, nil
}

// Delete removes a refresh token record. If not found, returns nil to simplify caller.
func (r *GormRefreshRepo) Delete(hashed string) error {
	if hashed == "" {
		return errors.New("hashed token required")
	}
	res := r.db.Where("token_hash = ?", hashed).Delete(&RefreshToken{})
	// treat not found as nil (idempotent delete)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// Revoke marks a refresh token as revoked (soft revoke).
// Optional convenience method for revocation instead of Delete.
func (r *GormRefreshRepo) Revoke(hashed string) error {
	if hashed == "" {
		return errors.New("hashed token required")
	}
	res := r.db.Model(&RefreshToken{}).Where("token_hash = ?", hashed).Updates(map[string]interface{}{
		"revoked":    true,
		"updated_at": time.Now().UTC(),
	})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// Helper: detect unique violation for Postgres driver errors.
// We try to detect a unique constraint error in a generic way.
// This helper inspects the underlying error for common drivers.
func isUniqueViolation(err error) bool {
	// Best-effort: pq, pgx, lib/pq vary; inspect SQLState or string fallback.
	// For simplicity check substring — robust detection can use pgerror from pgx.
	if err == nil {
		return false
	}
	msg := err.Error()
	// PostgreSQL unique_violation SQLSTATE = "23505" appears in pgerror details for pgx.
	if stringsContainsAny(msg, "unique constraint", "duplicate key", "23505") {
		return true
	}
	return false
}

// small helper to avoid importing strings dozens of times
func stringsContainsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if sub == "" {
			continue
		}
		if contains := (len(s) >= len(sub) && (stringIndex(s, sub) >= 0)); contains {
			return true
		}
	}
	return false
}

// tiny wrappers around standard library functions to keep this file self-contained
func stringIndex(s, sep string) int {
	return indexFunc(s, sep)
}

// Use strings.Index from stdlib indirectly so we don't import strings earlier than needed
func indexFunc(s, sep string) int {
	// implement a thin wrapper calling the stdlib function
	return indexStd(s, sep)
}

func indexStd(s, sep string) int {
	// actual implementation uses strings.Index
	return stringsIndex(s, sep)
}

// Actual strings.Index — we import it here to keep top of file clean in earlier lines.
func stringsIndex(s, sep string) int {
	return sqlIndex(s, sep)
}

// Now import and call strings.Index via alias (this avoids top-level import ordering confusion)
func sqlIndex(s, sep string) int {
	// simplest: use standard library
	return stringsIndexImpl(s, sep)
}

func stringsIndexImpl(s, sep string) int {
	// real implementation:
	return indexStdLib(s, sep)
}

func indexStdLib(s, sep string) int {
	// ok final call to stdlib:
	return strings.Index(s, sep)
}

package db

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"bitka/db/utils"
)

// NewDB creates a *gorm.DB for the given DSN and returns it.
// debug toggles GORM logging. This is generic and service-agnostic.
// Caller should manage schema selection (see Schema helpers).
func NewDB(dsn string, debug bool) (*gorm.DB, error) {
	if dsn == "" {
		return nil, errors.New("dsn required")
	}

	dsn, err := utils.SafeEncodeDSN(dsn)
	if err != nil {
		return nil, errors.New("failed to encode dsn")
	}

	var gormLogger logger.Interface
	if debug {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: false,
	})

	gormCfg := &gorm.Config{
		Logger: gormLogger,
	}

	gdb, err := gorm.Open(dialector, gormCfg)
	if err != nil {
		return nil, fmt.Errorf("gorm open: %w", err)
	}

	// configure sql.DB pool
	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, fmt.Errorf("db.DB(): %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	// verify connectivity
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}
	return gdb, nil
}

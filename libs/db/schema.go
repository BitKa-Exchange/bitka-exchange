package db

import (
	"fmt"

	"gorm.io/gorm"
)

// EnsureSchema creates the schema if not exists and sets search_path for the current connection.
// - schemaName: the schema the service will use (e.g. "auth").
// This function is service-agnostic; it doesn't know models.
func EnsureSchema(gdb *gorm.DB, schemaName string) error {
	if gdb == nil {
		return fmt.Errorf("db is nil")
	}
	// create schema if not exists
	if err := gdb.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName)).Error; err != nil {
		return err
	}
	// set search_path for session so subsequent operations (AutoMigrate, queries) use that schema
	if err := gdb.Exec(fmt.Sprintf("SET search_path TO %s", schemaName)).Error; err != nil {
		return err
	}
	return nil
}

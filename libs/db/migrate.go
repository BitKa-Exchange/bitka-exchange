package db

import "gorm.io/gorm"

// AutoMigrateModels runs AutoMigrate for the provided models.
// Keep migrations per service: each service passes only its models.
func AutoMigrateModels(gdb *gorm.DB, models ...interface{}) error {
	if gdb == nil {
		return nil
	}
	return gdb.AutoMigrate(models...)
}

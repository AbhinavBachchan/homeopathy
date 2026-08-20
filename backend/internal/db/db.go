package db

import (
	"homeopathy-platform/internal/config"
	"homeopathy-platform/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
}

// AutoMigrate creates/updates tables for Phase 1 scope.
// Consultation/subscription/corporate models get added here once Phase 2/3 start.
func AutoMigrate(conn *gorm.DB) error {
	return conn.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
	)
}

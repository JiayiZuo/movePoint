package database

import (
	"fmt"

	"movePoint/internal/models"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.User{},
		&models.ClimbingRecord{},
		&models.ClimbingAnalysis{},
	); err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	return nil
}

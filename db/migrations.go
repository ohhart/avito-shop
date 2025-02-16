package db

import (
	"fmt"

	"gorm.io/gorm"

	"avito-shop/internal/models"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Transaction{},
		&models.Inventory{},
		&models.Item{},
	)
	if err != nil {
		return fmt.Errorf("migration error: %w", err)
	}

	fmt.Println("Migration completed successfully!")
	return nil
}

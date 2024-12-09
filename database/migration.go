package database

import (
	"voucher_system/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.Voucher{},
		&models.Redeem{},
		&models.History{},
	)

	return err
}

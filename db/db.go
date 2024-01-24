package db

import (
	"web-service-ewallet/db/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	DB, err := gorm.Open(sqlite.Open("ewallet.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = DB.AutoMigrate(&models.Wallet{}, &models.Transaction{})
	if err != nil {
		panic("error while migrating")
	}

	return DB
}

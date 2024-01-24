package db

import (
	"web-service-ewallet/db/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	DB, err := gorm.Open(sqlite.Open("wallets.DB"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// migrate the schema
	err = DB.AutoMigrate(&models.Wallet{})

	if err != nil {
		panic("error while migrating")
	}

	return DB
}

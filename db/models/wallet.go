package models

import (
	"github.com/google/uuid"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Wallet struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Balance float64   `json:"balance"`
}

func CreateWallet(db *gorm.DB, wallet *Wallet) (*Wallet, error) {
	wallet.ID = uuid.New()
	err := db.Create(&wallet).Error
	return wallet, err
}

func GetWallet(db *gorm.DB, walletId string) (*Wallet, error) {
	var wallet Wallet
	err := db.First(&wallet, "id = ?", walletId).Error
	return &wallet, err
}

func GetAllWallets(db *gorm.DB) ([]Wallet, error) {
	var wallets []Wallet
	err := db.Find(&wallets).Error
	return wallets, err
}

func UpdateWallet(db *gorm.DB, wallet *Wallet) (*Wallet, error) {
	err := db.Save(wallet).Error
	return wallet, err
}

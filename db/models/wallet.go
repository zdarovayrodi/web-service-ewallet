package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

func GetAllWallets(c *gin.Context, db *gorm.DB) {
	var wallets []Wallet
	db.Find(&wallets)
	c.IndentedJSON(http.StatusOK, gin.H{"data": wallets})
}

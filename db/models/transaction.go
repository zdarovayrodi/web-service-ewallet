package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Amount            float64   `json:"amount"`
	SenderWalletID    uuid.UUID `json:"sender_wallet_id" gorm:"type:uuid"`
	RecipientWalletID uuid.UUID `json:"recipient_wallet_id" gorm:"type:uuid"`
}

func CreateTransaction(db *gorm.DB, senderID uuid.UUID, recipientID uuid.UUID, amount float64) (*Transaction, error) {
	transaction := &Transaction{
		ID:                uuid.New(),
		Amount:            amount,
		SenderWalletID:    senderID,
		RecipientWalletID: recipientID,
	}

	err := db.Create(&transaction).Error
	return transaction, err
}

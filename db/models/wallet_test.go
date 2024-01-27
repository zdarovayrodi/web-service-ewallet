package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	// in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database for testing")
	}

	err = db.AutoMigrate(&Wallet{})
	if err != nil {
		panic("error while migrating")
	}

	return db
}

func TestCreateWallet(t *testing.T) {
	db := setupTestDB()

	wallet := &Wallet{
		Balance: 100.0,
	}

	createdWallet, err := CreateWallet(db, wallet)
	defer db.Delete(createdWallet)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, createdWallet.ID)
	assert.Equal(t, wallet.Balance, createdWallet.Balance)
}

func TestGetWallet(t *testing.T) {
	db := setupTestDB()

	initialWallet := &Wallet{
		Balance: 200.0,
	}

	createdWallet, err := CreateWallet(db, initialWallet)
	assert.NoError(t, err)
	defer db.Delete(createdWallet)

	walletID := createdWallet.ID.String()

	retrievedWallet, err := GetWallet(db, walletID)
	assert.NoError(t, err)
	assert.Equal(t, createdWallet.ID, retrievedWallet.ID)
	assert.Equal(t, createdWallet.Balance, retrievedWallet.Balance)
}

func TestUpdateWallet(t *testing.T) {
	db := setupTestDB()

	initialWallet := &Wallet{
		Balance: 300.0,
	}

	createdWallet, err := CreateWallet(db, initialWallet)
	assert.NoError(t, err)
	defer db.Delete(createdWallet)

	createdWallet.Balance = 400.0
	updatedWallet, err := UpdateWallet(db, createdWallet)
	assert.NoError(t, err)
	assert.Equal(t, createdWallet.ID, updatedWallet.ID)
	assert.Equal(t, createdWallet.Balance, updatedWallet.Balance)

	retrievedWallet, err := GetWallet(db, createdWallet.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, createdWallet.Balance, retrievedWallet.Balance)
}

func TestGetAllWallets(t *testing.T) {
	db := setupTestDB()

	wallets := []*Wallet{
		{Balance: 500.0},
		{Balance: 600.0},
		{Balance: 700.0},
	}

	for _, wallet := range wallets {
		_, err := CreateWallet(db, wallet)
		assert.NoError(t, err)
		defer func(walletID uuid.UUID) {
			// delete the wallet after the test
			err := db.Where("id = ?", walletID).Delete(&Wallet{}).Error
			assert.NoError(t, err)
		}(wallet.ID)
	}

	allWallets, err := GetAllWallets(db)
	assert.NoError(t, err)

	assert.Equal(t, len(wallets), len(allWallets))

	for i, wallet := range wallets {
		assert.Equal(t, wallet.Balance, allWallets[i].Balance)
	}
}

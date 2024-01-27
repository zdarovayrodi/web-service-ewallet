// models/transaction_test.go

package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	db := setupTestDB()

	senderWallet := &Wallet{Balance: 1000.0}
	recipientWallet := &Wallet{Balance: 500.0}

	_, err := CreateWallet(db, senderWallet)
	assert.NoError(t, err)
	defer db.Delete(senderWallet)

	_, err = CreateWallet(db, recipientWallet)
	assert.NoError(t, err)
	defer db.Delete(recipientWallet)

	amount := 200.0
	createdTransaction, err := CreateTransaction(db, senderWallet.ID, recipientWallet.ID, amount)
	assert.NoError(t, err)
	defer db.Delete(createdTransaction)

	assert.NotEqual(t, uuid.Nil, createdTransaction.ID)
	assert.Equal(t, amount, createdTransaction.Amount)
	assert.Equal(t, senderWallet.ID, createdTransaction.SenderWalletID)
	assert.Equal(t, recipientWallet.ID, createdTransaction.RecipientWalletID)
}

func TestGetIncomingTransactions(t *testing.T) {
	db := setupTestDB()

	senderWallet := &Wallet{Balance: 1000.0}
	recipientWallet := &Wallet{Balance: 500.0}

	_, err := CreateWallet(db, senderWallet)
	assert.NoError(t, err)
	defer db.Delete(senderWallet)

	_, err = CreateWallet(db, recipientWallet)
	assert.NoError(t, err)
	defer db.Delete(recipientWallet)

	// incoming transactions
	incomingTransaction1, err := CreateTransaction(db, senderWallet.ID, recipientWallet.ID, 200.0)
	assert.NoError(t, err)
	defer db.Delete(incomingTransaction1)

	incomingTransaction2, err := CreateTransaction(db, senderWallet.ID, recipientWallet.ID, 300.0)
	assert.NoError(t, err)
	defer db.Delete(incomingTransaction2)

	// outgoing transaction
	outgoingTransaction, err := CreateTransaction(db, recipientWallet.ID, senderWallet.ID, 50.0)
	assert.NoError(t, err)
	defer db.Delete(outgoingTransaction)

	incomingTransactions, err := GetIncomingTransactions(db, recipientWallet.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(incomingTransactions))

	assert.Equal(t, incomingTransaction1.ID, incomingTransactions[0].ID)
	assert.Equal(t, incomingTransaction2.ID, incomingTransactions[1].ID)
}

func TestGetOutgoingTransactions(t *testing.T) {
	db := setupTestDB()

	senderWallet := &Wallet{Balance: 1000.0}
	recipientWallet := &Wallet{Balance: 500.0}

	_, err := CreateWallet(db, senderWallet)
	assert.NoError(t, err)
	defer db.Delete(senderWallet)

	_, err = CreateWallet(db, recipientWallet)
	assert.NoError(t, err)
	defer db.Delete(recipientWallet)

	// outgoing transactions
	outgoingTransaction1, err := CreateTransaction(db, senderWallet.ID, recipientWallet.ID, 200.0)
	assert.NoError(t, err)
	defer db.Delete(outgoingTransaction1)

	outgoingTransaction2, err := CreateTransaction(db, senderWallet.ID, recipientWallet.ID, 300.0)
	assert.NoError(t, err)
	defer db.Delete(outgoingTransaction2)

	// incoming transaction
	incomingTransaction, err := CreateTransaction(db, recipientWallet.ID, senderWallet.ID, 50.0)
	assert.NoError(t, err)
	defer db.Delete(incomingTransaction)

	outgoingTransactions, err := GetOutgoingTransactions(db, senderWallet.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(outgoingTransactions))

	assert.Equal(t, outgoingTransaction1.ID, outgoingTransactions[0].ID)
	assert.Equal(t, outgoingTransaction2.ID, outgoingTransactions[1].ID)
}

package main

import (
	"net/http"

	"web-service-ewallet/db"
	"web-service-ewallet/db/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

// GET /api/v1/wallets
func getWallets(context *gin.Context) {
	models.GetAllWallets(context, DB)
}

// GET /api/v1/wallets/{id}
func getWalletById(context *gin.Context) {
	id := context.Param("id")
	wallet, err := models.GetWallet(DB, id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Wallet not found"})
		return
	}

	context.JSON(http.StatusOK, wallet)
}

// POST /api/v1/wallets
func postWallet(context *gin.Context) {
	var newWallet models.Wallet

	if context.Request.ContentLength == 0 {
		newWallet.Balance = 100
	} else {
		// in case "balance" field is frovided
		if err := context.BindJSON(&newWallet); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
	}

	createdWallet, err := models.CreateWallet(DB, &newWallet)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	context.JSON(http.StatusCreated, createdWallet)
}

// Перевод средств с одного кошелька на другой
// Эндпоинт - POST /api/v1/wallet/{walletId}/send
// Параметры запроса:
// walletId – строковый ID кошелька, указан в пути запроса
// JSON-объект в теле запроса с параметрами:
// to – ID кошелька, куда нужно перевести деньги
// amount – сумма перевода
// Статус ответа 200 если перевод успешен
// Статус ответа 404 если исходящий кошелек не найден
// Статус ответа 400 если целевой кошелек не найден или на исходящем нет нужной суммы

// POST /api/v1/wallet/{walletId}/send
func transferFunds(context *gin.Context) {
	walletID := context.Param("walletID")

	fromWallet, err := models.GetWallet(DB, walletID)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Source wallet not found"})
		return
	}

	var transferFundsRequest struct {
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	}

	if err := context.BindJSON(&transferFundsRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	toWallet, err := models.GetWallet(DB, transferFundsRequest.To)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Destination wallet not found"})
		return
	}

	if fromWallet.Balance < transferFundsRequest.Amount {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
		return
	}

	_, err = models.CreateTransaction(DB, fromWallet.ID, toWallet.ID, transferFundsRequest.Amount)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	fromWallet.Balance -= transferFundsRequest.Amount
	toWallet.Balance += transferFundsRequest.Amount

	_, err = models.UpdateWallet(DB, fromWallet)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update source wallet"})
		return
	}

	_, err = models.UpdateWallet(DB, toWallet)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update destination wallet"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "Transfer successful"})
}

// Получение историй входящих и исходящих транзакций
// Эндпоинт – GET /api/v1/wallet/{walletId}/history
// Параметры запроса:
// walletId – строковый ID кошелька, указан в пути запроса
// Ответ с статусом 200 если кошелек найден. Ответ должен содержать в теле массив JSON-объектов с входящими и исходящими транзакциями кошелька. Каждый объект содержит параметры:
// time – дата и время перевода в формате RFC 3339
// from – ID исходящего кошелька
// to – ID входящего кошелька
// amount – сумма перевода. Дробное число
// Статус ответа 404 если указанный кошелек не найден

// GET /api/v1/wallet/{walletId}/history
func getTransactionHistory(context *gin.Context) {
	walletID := context.Param("walletID")

	wallet, err := models.GetWallet(DB, walletID)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
	}

	var incomingTransactions []models.Transaction
	var outgoingTransactions []models.Transaction

	incomingTransactions, err = models.GetIncomingTransactions(DB, wallet.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve incoming transactions"})
		return
	}

	outgoingTransactions, err = models.GetOutgoingTransactions(DB, wallet.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve outgoing transactions"})
		return
	}

	transactionHistory := append(incomingTransactions, outgoingTransactions...)

	context.IndentedJSON(http.StatusOK, gin.H{"transactionHistory": transactionHistory})
}

func main() {
	//db
	DB = db.InitDB()

	// api
	router := gin.Default()
	router.GET("/api/v1/wallets", getWallets)
	router.POST("/api/v1/wallets", postWallet)
	router.GET("/api/v1/wallets/:walletID/history", getTransactionHistory)
	router.POST("/api/v1/wallets/:walletID/send", transferFunds)
	router.GET("/api/v1/wallets/:walletID", getWalletById)

	router.Run("localhost:8080")
}

package main

import (
	"net/http"

	"web-service-ewallet/db/models"

	"github.com/gin-gonic/gin"
)

var wallets = []models.Wallet{
	{ID: "1", Balance: 56.01},
	{ID: "2", Balance: 156.02},
	{ID: "3", Balance: 1999887799956.03},
}

// /api/v1/wallets (get all)
func getWallets(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, wallets)
}

// /api/v1/wallets/{id}
func getWalletById(context *gin.Context) {
	id := context.Param("id")

	for _, wallet := range wallets {
		if wallet.ID == id {
			context.IndentedJSON(http.StatusOK, wallet)
			return
		}
	}
	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "wallet not found"})
}

// /api/v1/wallets (post)
func postWallet(context *gin.Context) {
	var newWallet models.Wallet

	if err := context.BindJSON(&newWallet); err != nil {
		return
	}

	wallets = append(wallets, newWallet)
	context.IndentedJSON(http.StatusCreated, newWallet)
}

func main() {
	router := gin.Default()
	router.GET("/api/v1/wallets", getWallets)
	router.GET("/api/v1/wallets/:id", getWalletById)
	router.POST("/api/v1/wallets", postWallet)

	router.Run("localhost:8080")
}

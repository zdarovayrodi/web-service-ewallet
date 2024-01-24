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

func main() {
	//db
	DB = db.InitDB()

	// api
	router := gin.Default()
	router.GET("/api/v1/wallets", getWallets)
	router.GET("/api/v1/wallets/:id", getWalletById)
	router.POST("/api/v1/wallets", postWallet)

	router.Run("localhost:8080")
}

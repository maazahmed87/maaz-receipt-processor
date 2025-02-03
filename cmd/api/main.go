package main

import (
	"log"
	"net/http"
	"receipt-processor/internal/api/handler"
	"receipt-processor/internal/domain/service"
	"receipt-processor/internal/storage/memory"
	"time"

	"github.com/gin-gonic/gin"
)

func setupHealthCheck(router *gin.Engine) {
	router.GET("receipts/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now(),
		})
	})
}

func main() {
	// initialize the dependencies
	storage := memory.NewMemoryStorage()
	receiptService := service.NewReceiptService(storage)
	receiptHandler := handler.NewReceiptHandler(receiptService)

	// defined routes
	// create a default gin router
	router := gin.Default()
	setupHealthCheck(router)
	router.POST("/receipts/process", receiptHandler.ProcessReceipt)
	router.GET("/receipts/:id/points", receiptHandler.GetPoints)

	// starting the server
	log.Println("Server starting on port 8080-----")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

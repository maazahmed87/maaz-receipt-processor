package main

import (
	"log"
	"receipt-processor/internal/api/handler"
	"receipt-processor/internal/domain/service"
	"receipt-processor/internal/storage/memory"

	"github.com/gin-gonic/gin"
)

func main() {
	// initialize the dependencies
	storage := memory.NewMemoryStorage()
	receiptService := service.NewReceiptService(storage)
	receiptHandler := handler.NewReceiptHandler(receiptService)

	// defined routes
	// create a default gin router
	router := gin.Default()
	router.GET("/receipts/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})
	router.POST("/receipts/process", receiptHandler.ProcessReceipt)
	router.GET("/receipts/:id/points", receiptHandler.GetPoints)

	// starting the server
	log.Println("Server starting on port 8080-----")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

package main

import (
	"log"
	"receipt-processor/internal/api/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	// create a default gin router
	router := gin.Default()

	// initialize the receipt handler
	receiptHandler := handler.NewReceiptHandler()

	// defined routes
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

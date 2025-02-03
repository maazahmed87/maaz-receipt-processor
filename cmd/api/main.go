package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"receipt-processor/internal/api/handler"
	"receipt-processor/internal/domain/service"
	"receipt-processor/internal/storage/memory"
	"syscall"
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

func setupRouter(receiptHandler *handler.ReceiptHandler) *gin.Engine {
	router := gin.Default()
	setupHealthCheck(router)
	router.POST("/receipts/process", receiptHandler.ProcessReceipt)
	router.GET("/receipts/:id/points", receiptHandler.GetPoints)
	return router
}

func main() {
	// initialize the dependencies
	storage := memory.NewMemoryStorage()
	receiptService := service.NewReceiptService(storage)
	receiptHandler := handler.NewReceiptHandler(receiptService)

	// setting up router
	router := setupRouter(receiptHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Starting server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// using the interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Giving outstanding requests 5 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting...")
}

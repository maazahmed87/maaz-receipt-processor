package handler

import (
	"net/http"
	"receipt-processor/internal/domain/model"

	"github.com/gin-gonic/gin"
)

// ReceiptHandler hanldes HTTP requests for receipts
type ReceiptHandler struct {
}

// NewReceiptHanlder creates a new receipt handler
func NewReceiptHandler() *ReceiptHandler {
	return &ReceiptHandler{}
}

func (h *ReceiptHandler) ProcessReceipt(c *gin.Context) {
	var receipt model.Receipt

	if err := c.ShouldBindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid receipt format",
		})
		return
	}

	//returning dummy-id for now
	response := model.ProcessResponse{
		ID: "dummy-id",
	}

	c.JSON(http.StatusOK, response)
}

func (h *ReceiptHandler) GetPoints(c *gin.Context) {
	// id := c.Param("id")

	// returning dummy points for now
	response := model.PointsResponse{
		Points: 100,
	}

	c.JSON(http.StatusOK, response)
}

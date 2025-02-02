package handler

import (
	"net/http"
	"receipt-processor/internal/domain/model"
	"receipt-processor/internal/domain/service"
	"receipt-processor/pkg/validator"

	"github.com/gin-gonic/gin"
)

// ReceiptHandler hanldes HTTP requests for receipts
type ReceiptHandler struct {
	service *service.ReceiptService
}

// NewReceiptHanlder creates a new receipt handler
func NewReceiptHandler(service *service.ReceiptService) *ReceiptHandler {
	return &ReceiptHandler{
		service: service,
	}
}

func (h *ReceiptHandler) ProcessReceipt(c *gin.Context) {
	var receipt model.Receipt

	if err := c.ShouldBindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid receipt format",
		})
		return
	}

	id, err := h.service.ProcessReceipt(receipt)

	if err != nil {
		if validationErr, ok := err.(*validator.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr.Error(),
			})
			return
		}
		c.JSON((http.StatusInternalServerError), gin.H{"error": "Failed to process receipt"})
		return
	}

	c.JSON(http.StatusOK, model.ProcessResponse{ID: id})
}

func (h *ReceiptHandler) GetPoints(c *gin.Context) {
	id := c.Param("id")

	points, err := h.service.GetPoints(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
		return
	}

	c.JSON(http.StatusOK, model.PointsResponse{Points: points})
}

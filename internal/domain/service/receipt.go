package service

import (
	"math"
	"receipt-processor/internal/domain/model"
	"receipt-processor/internal/storage/memory"
	"receipt-processor/pkg/validator"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ReceiptService hanles the business logic for receipt processing
type ReceiptService struct {
	storage   memory.ReceiptStorage
	validator *validator.ReceiptValidator
}

// NewReceiptService creates a new ReceiptService instance
func NewReceiptService(storage memory.ReceiptStorage) *ReceiptService {
	return &ReceiptService{
		storage:   storage,
		validator: validator.NewReceiptValidator(),
	}
}

// ProcessReceipt calculates the points for a receipt and stores the receipt
func (s *ReceiptService) ProcessReceipt(receipt model.Receipt) (string, error) {
	id := uuid.New().String()

	if err := s.validator.Validate(receipt); err != nil {
		return "", err
	}

	points := s.calculatePoints(receipt)

	if err := s.storage.SavePoints(id, points); err != nil {
		return "", err
	}
	return id, nil
}

// GetPoints returns the points for a receipt id
func (s *ReceiptService) GetPoints(id string) (int, error) {
	return s.storage.GetPoints(id)
}

func (s *ReceiptService) calculatePoints(receipt model.Receipt) int {
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name.
	alphanumericRegex := regexp.MustCompile("[a-zA-Z0-9]")
	points += len(alphanumericRegex.FindAllString(receipt.Retailer, -1))

	// Rule 2: 50 points if the total is a round dollar amount with no cents.
	if s.isRoundDollarAmount(receipt.Total) {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25.
	if s.isMultipleOf25Cents(receipt.Total) {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt.
	points += (len(receipt.Items) / 2) * 5

	// Rule 5: If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2
	// and round up to the nearest integer. The result is the number of points earned.
	for _, item := range receipt.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd.
	if s.isOddDay(receipt.PurchaseDate) {
		points += 6
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	if s.isAfternoonPurchase(receipt.PurchaseTime) {
		points += 10
	}

	return points
}

// helper methods for calculating points using the rules

func (s *ReceiptService) isRoundDollarAmount(total string) bool {
	if value, err := strconv.ParseFloat(total, 64); err == nil {
		return math.Mod(value, 1.0) == 0
	}
	return false
}

func (s *ReceiptService) isMultipleOf25Cents(total string) bool {
	if value, err := strconv.ParseFloat(total, 64); err == nil {
		return math.Mod(value*100, 25) == 0
	}
	return false
}

func (s *ReceiptService) isOddDay(date string) bool {
	if purchaseDate, err := time.Parse("2006-01-02", date); err == nil {
		return purchaseDate.Day()%2 == 1
	}
	return false
}

func (s *ReceiptService) isAfternoonPurchase(timeStr string) bool {
	if purchaseTime, err := time.Parse("15:04", timeStr); err == nil {
		hour := purchaseTime.Hour()
		return hour >= 14 && hour < 16
	}
	return false
}

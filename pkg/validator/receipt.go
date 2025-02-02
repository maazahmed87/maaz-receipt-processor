package validator

import (
	"fmt"
	"receipt-processor/internal/domain/model"
	"regexp"
	"strings"
	"time"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type ReceiptValidator struct{}

func NewReceiptValidator() *ReceiptValidator {
	return &ReceiptValidator{}
}

func (v *ReceiptValidator) Validate(receipt model.Receipt) error {
	// validate retailer
	if strings.TrimSpace(receipt.Retailer) == "" {
		return &ValidationError{
			Field:   "retailer",
			Message: "retailer is required",
		}
	}

	// Validate purchase date
	if err := v.validateDate(receipt.PurchaseDate); err != nil {
		return &ValidationError{
			Field:   "purchaseDate",
			Message: err.Error(),
		}
	}

	// validate purchase time
	if err := v.validateTime(receipt.PurchaseTime); err != nil {
		return &ValidationError{
			Field:   "purchaseTime",
			Message: err.Error(),
		}
	}

	// validate ietms count
	if len(receipt.Items) == 0 {
		return &ValidationError{
			Field:   "items",
			Message: "At least one item is required",
		}
	}

	// validating each item
	for i, item := range receipt.Items {
		if err := v.validateItem(item, i); err != nil {
			return err
		}
	}

	// Validating total
	if err := v.validateTotal(receipt.Total); err != nil {
		return &ValidationError{
			Field:   "total",
			Message: err.Error(),
		}
	}

	return nil
}

func (v *ReceiptValidator) validateDate(date string) error {
	if date == "" {
		return fmt.Errorf("purchase date is required")
	}

	_, err := time.Parse("2006-01-02", date)

	if err != nil {
		return fmt.Errorf("invalid date format, expecter YYYY-MM-DD")
	}

	return nil
}

func (v *ReceiptValidator) validateTime(timeStr string) error {
	if timeStr == "" {
		return fmt.Errorf("purchase time is required")
	}
	_, err := time.Parse("15:04", timeStr)

	if err != nil {
		return fmt.Errorf("invalid time format, expected HH:MM")
	}
	return nil
}

func (v *ReceiptValidator) validateItem(item model.Item, index int) error {
	if strings.TrimSpace(item.ShortDescription) == "" {
		return &ValidationError{
			Field:   fmt.Sprintf("items[%d].shortDescription", index),
			Message: "item description is required",
		}
	}
	if err := v.validatePrice(item.Price); err != nil {
		return &ValidationError{
			Field:   fmt.Sprintf("items[%d].price", index),
			Message: err.Error(),
		}
	}
	return nil
}

func (v *ReceiptValidator) validatePrice(price string) error {
	if price == "" {
		return fmt.Errorf("price is required")
	}

	matched, err := regexp.MatchString(`^\d+\.\d{2}$`, price)
	if err != nil || !matched {
		return fmt.Errorf("invalid price format")
	}
	return nil
}

func (v *ReceiptValidator) validateTotal(total string) error {
	if total == "" {
		return fmt.Errorf("total is required")
	}

	matched, err := regexp.MatchString(`^\d+\.\d{2}$`, total)

	if err != nil || !matched {
		return fmt.Errorf("invalid total format")
	}
	return nil
}

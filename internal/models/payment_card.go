package models

import (
	"errors"
	"regexp"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentCard struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID         primitive.ObjectID `json:"user_id" bson:"user_id"`
	CardHolderName string             `json:"card_holder_name" bson:"card_holder_name"`
	CardNumber     string             `json:"card_number" bson:"card_number"`
	ExpiryDate     string             `json:"expiry_date" bson:"expiry_date"`
	CVV            string             `json:"cvv" bson:"cvv"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
}

type PaymentCardCreate struct {
	CardHolderName string `json:"card_holder_name" binding:"required"`
	CardNumber     string `json:"card_number" binding:"required"`
	ExpiryDate     string `json:"expiry_date" binding:"required"`
	CVV            string `json:"cvv" binding:"required"`
}

func ValidateCardNumber(cardNumber string) error {
	if len(cardNumber) != 16 {
		return errors.New("card number must be exactly 16 digits")
	}

	matched, _ := regexp.MatchString(`^\d{16}$`, cardNumber)
	if !matched {
		return errors.New("card number must contain only digits")
	}

	return nil
}

func ValidateExpiryDate(expiryDate string) error {
	matched, _ := regexp.MatchString(`^\d{2}/\d{2}$`, expiryDate)
	if !matched {
		return errors.New("expiry date must be in MM/YY format")
	}

	month, err := strconv.Atoi(expiryDate[0:2])
	if err != nil || month < 1 || month > 12 {
		return errors.New("invalid month in expiry date (must be 01-12)")
	}

	year, err := strconv.Atoi(expiryDate[3:5])
	if err != nil {
		return errors.New("invalid year in expiry date")
	}

	now := time.Now()
	currentYear := now.Year() % 100
	currentMonth := int(now.Month())

	if year < currentYear || (year == currentYear && month < currentMonth) {
		return errors.New("card has expired")
	}

	return nil
}

func ValidateCVV(cvv string) error {
	if len(cvv) < 3 || len(cvv) > 4 {
		return errors.New("CVV must be 3 or 4 digits")
	}

	matched, _ := regexp.MatchString(`^\d{3,4}$`, cvv)
	if !matched {
		return errors.New("CVV must contain only digits")
	}

	return nil
}

func (pc *PaymentCardCreate) Validate() error {
	if pc.CardHolderName == "" {
		return errors.New("card holder name is required")
	}

	if err := ValidateCardNumber(pc.CardNumber); err != nil {
		return err
	}

	if err := ValidateExpiryDate(pc.ExpiryDate); err != nil {
		return err
	}

	if err := ValidateCVV(pc.CVV); err != nil {
		return err
	}

	return nil
}

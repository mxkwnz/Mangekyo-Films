package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "PENDING"
	PaymentCompleted PaymentStatus = "COMPLETED"
	PaymentFailed    PaymentStatus = "FAILED"
	PaymentRefunded  PaymentStatus = "REFUNDED"
)

const (
	MaxPaymentAmount = 100000.0
)

type Payment struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID          primitive.ObjectID `json:"user_id" bson:"user_id"`
	PaymentCardID   primitive.ObjectID `json:"payment_card_id" bson:"payment_card_id"`
	TransactionCode string             `json:"transaction_code" bson:"transaction_code"`
	Amount          float64            `json:"amount" bson:"amount"`
	Status          PaymentStatus      `json:"status" bson:"status"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
}

type PaymentCreate struct {
	PaymentCardID primitive.ObjectID `json:"payment_card_id" binding:"required"`
	Amount        float64            `json:"amount" binding:"required,gt=0"`
}

func ValidateAmount(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	if amount > MaxPaymentAmount {
		return errors.New("amount exceeds maximum limit of 100,000 tenge")
	}

	return nil
}

func (pc *PaymentCreate) Validate() error {
	if pc.PaymentCardID.IsZero() {
		return errors.New("payment card ID is required")
	}

	if err := ValidateAmount(pc.Amount); err != nil {
		return err
	}

	return nil
}

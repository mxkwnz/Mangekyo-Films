package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentCode struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Code      string             `json:"code" bson:"code"`
	Amount    float64            `json:"amount" bson:"amount"`
	IsUsed    bool               `json:"is_used" bson:"is_used"`
	UsedBy    primitive.ObjectID `json:"used_by" bson:"used_by,omitempty"`
	UsedAt    time.Time          `json:"used_at" bson:"used_at,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type RedeemCodeRequest struct {
	Code string `json:"code" binding:"required"`
}

func (pc *PaymentCode) Validate() error {
	if pc.Code == "" {
		return errors.New("code is required")
	}
	if pc.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	return nil
}

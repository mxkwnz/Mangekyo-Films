package repositories

import (
	"cinema-system/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentCodeRepository struct {
	collection *mongo.Collection
}

func NewPaymentCodeRepository(db *mongo.Database) *PaymentCodeRepository {
	return &PaymentCodeRepository{
		collection: db.Collection("payment_codes"),
	}
}

func (r *PaymentCodeRepository) FindByCode(ctx context.Context, code string) (*models.PaymentCode, error) {
	var pc models.PaymentCode
	err := r.collection.FindOne(ctx, bson.M{"code": code}).Decode(&pc)
	if err != nil {
		return nil, err
	}
	return &pc, nil
}

func (r *PaymentCodeRepository) MarkAsUsed(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"is_used": true,
				"used_by": userID,
				"used_at": time.Now(),
			},
		},
	)
	return err
}

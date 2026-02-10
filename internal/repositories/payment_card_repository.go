package repositories

import (
	"cinema-system/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentCardRepository struct {
	collection *mongo.Collection
}

func NewPaymentCardRepository(db *mongo.Database) *PaymentCardRepository {
	return &PaymentCardRepository{
		collection: db.Collection("payment_cards"),
	}
}

func (r *PaymentCardRepository) Create(ctx context.Context, card *models.PaymentCard) error {
	result, err := r.collection.InsertOne(ctx, card)
	if err != nil {
		return err
	}
	card.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *PaymentCardRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.PaymentCard, error) {
	var card models.PaymentCard
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&card)
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *PaymentCardRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]models.PaymentCard, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cards []models.PaymentCard
	if err = cursor.All(ctx, &cards); err != nil {
		return nil, err
	}
	return cards, nil
}

func (r *PaymentCardRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *PaymentCardRepository) Update(ctx context.Context, id primitive.ObjectID, card *models.PaymentCard) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": card},
	)
	return err
}
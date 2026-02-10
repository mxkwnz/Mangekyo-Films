package repositories

import (
	"cinema-system/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentRepository struct {
	collection *mongo.Collection
}

func NewPaymentRepository(db *mongo.Database) *PaymentRepository {
	return &PaymentRepository{
		collection: db.Collection("payments"),
	}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	result, err := r.collection.InsertOne(ctx, payment)
	if err != nil {
		return err
	}
	payment.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *PaymentRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Payment, error) {
	var payment models.Payment
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&payment)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]models.Payment, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []models.Payment
	if err = cursor.All(ctx, &payments); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepository) FindByTransactionCode(ctx context.Context, code string) (*models.Payment, error) {
	var payment models.Payment
	err := r.collection.FindOne(ctx, bson.M{"transaction_code": code}).Decode(&payment)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status models.PaymentStatus) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"status": status}},
	)
	return err
}

func (r *PaymentRepository) GetAll(ctx context.Context) ([]models.Payment, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []models.Payment
	if err = cursor.All(ctx, &payments); err != nil {
		return nil, err
	}
	return payments, nil
}
package repositories

import (
	"cinema-system/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HallRepository struct {
	collection *mongo.Collection
}

func NewHallRepository(db *mongo.Database) *HallRepository {
	return &HallRepository{
		collection: db.Collection("halls"),
	}
}

func (r *HallRepository) Create(ctx context.Context, hall *models.Hall) error {
	result, err := r.collection.InsertOne(ctx, hall)
	if err != nil {
		return err
	}
	hall.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *HallRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Hall, error) {
	var hall models.Hall
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&hall)
	if err != nil {
		return nil, err
	}
	return &hall, nil
}

func (r *HallRepository) GetAll(ctx context.Context) ([]models.Hall, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var halls []models.Hall
	if err = cursor.All(ctx, &halls); err != nil {
		return nil, err
	}
	return halls, nil
}

func (r *HallRepository) Update(ctx context.Context, id primitive.ObjectID, hall *models.Hall) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": hall},
	)
	return err
}

func (r *HallRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

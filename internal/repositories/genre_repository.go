package repositories

import (
	"cinema-system/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GenreRepository struct {
	collection *mongo.Collection
}

func NewGenreRepository(db *mongo.Database) *GenreRepository {
	return &GenreRepository{
		collection: db.Collection("genres"),
	}
}

func (r *GenreRepository) Create(ctx context.Context, genre *models.Genre) error {
	result, err := r.collection.InsertOne(ctx, genre)
	if err != nil {
		return err
	}
	genre.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *GenreRepository) GetAll(ctx context.Context) ([]models.Genre, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var genres []models.Genre
	if err = cursor.All(ctx, &genres); err != nil {
		return nil, err
	}
	return genres, nil
}

func (r *GenreRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Genre, error) {
	var genre models.Genre
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&genre)
	if err != nil {
		return nil, err
	}
	return &genre, nil
}

func (r *GenreRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

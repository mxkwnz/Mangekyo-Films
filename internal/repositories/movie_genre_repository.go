package repositories

import (
	"cinema-system/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MovieGenreRepository struct {
	collection *mongo.Collection
}

func NewMovieGenreRepository(db *mongo.Database) *MovieGenreRepository {
	return &MovieGenreRepository{
		collection: db.Collection("movie_genres"),
	}
}

func (r *MovieGenreRepository) Create(ctx context.Context, movieGenre *models.MovieGenre) error {
	result, err := r.collection.InsertOne(ctx, movieGenre)
	if err != nil {
		return err
	}
	movieGenre.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *MovieGenreRepository) GetGenresByMovieID(ctx context.Context, movieID primitive.ObjectID) ([]models.MovieGenre, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"movie_id": movieID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movieGenres []models.MovieGenre
	if err = cursor.All(ctx, &movieGenres); err != nil {
		return nil, err
	}
	return movieGenres, nil
}

func (r *MovieGenreRepository) GetMoviesByGenreID(ctx context.Context, genreID primitive.ObjectID) ([]models.MovieGenre, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"genre_id": genreID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movieGenres []models.MovieGenre
	if err = cursor.All(ctx, &movieGenres); err != nil {
		return nil, err
	}
	return movieGenres, nil
}

func (r *MovieGenreRepository) DeleteByMovieID(ctx context.Context, movieID primitive.ObjectID) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"movie_id": movieID})
	return err
}

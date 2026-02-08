package repositories

import (
	"cinema-system/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MovieRepository struct {
	collection *mongo.Collection
}

func NewMovieRepository(db *mongo.Database) *MovieRepository {
	return &MovieRepository{
		collection: db.Collection("movies"),
	}
}

func (r *MovieRepository) Create(ctx context.Context, movie *models.Movie) error {
	result, err := r.collection.InsertOne(ctx, movie)
	if err != nil {
		return err
	}
	movie.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *MovieRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Movie, error) {
	var movie models.Movie
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&movie)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *MovieRepository) GetAll(ctx context.Context) ([]models.Movie, error) {
	findOptions := options.Find().SetSort(bson.D{{Key: "rating", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movies []models.Movie
	if err = cursor.All(ctx, &movies); err != nil {
		return nil, err
	}
	return movies, nil
}

func (r *MovieRepository) Update(ctx context.Context, id primitive.ObjectID, movie *models.Movie) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": movie},
	)
	return err
}

func (r *MovieRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *MovieRepository) UpdateRating(ctx context.Context, movieID primitive.ObjectID, newRating float64) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": movieID},
		bson.M{"$set": bson.M{"rating": newRating}},
	)
	return err
}

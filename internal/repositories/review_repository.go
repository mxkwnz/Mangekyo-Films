package repositories

import (
	"cinema-system/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewRepository struct {
	collection *mongo.Collection
}

func NewReviewRepository(db *mongo.Database) *ReviewRepository {
	return &ReviewRepository{
		collection: db.Collection("reviews"),
	}
}

func (r *ReviewRepository) Create(ctx context.Context, review *models.Review) error {
	result, err := r.collection.InsertOne(ctx, review)
	if err != nil {
		return err
	}
	review.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ReviewRepository) GetByMovie(ctx context.Context, movieID primitive.ObjectID) ([]models.Review, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"movie_id": movieID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reviews []models.Review
	if err = cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *ReviewRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Review, error) {
	var review models.Review
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&review)
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *ReviewRepository) GetAverageRating(ctx context.Context, movieID primitive.ObjectID) (float64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"movie_id": movieID}}},
		{{Key: "$group", Value: bson.M{
			"_id":       nil,
			"avgRating": bson.M{"$avg": "$rating"},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		AvgRating float64 `bson:"avgRating"`
	}
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}

	return result[0].AvgRating, nil
}

func (r *ReviewRepository) GetByUser(ctx context.Context, userID primitive.ObjectID) ([]models.Review, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"user_id": userID},
			{"user_id": userID.Hex()},
		},
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reviews []models.Review
	if err = cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *ReviewRepository) CheckUserReview(ctx context.Context, userID, movieID primitive.ObjectID) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"user_id":  userID,
		"movie_id": movieID,
	})
	return count > 0, err
}

func (r *ReviewRepository) Update(ctx context.Context, review *models.Review) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": review.ID},
		bson.M{"$set": review},
	)
	return err
}

func (r *ReviewRepository) UpdateReviewerName(ctx context.Context, userID primitive.ObjectID, newName string) error {
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{
			"$or": []bson.M{
				{"user_id": userID},
				{"user_id": userID.Hex()},
			},
		},
		bson.M{"$set": bson.M{"user_name": newName}},
	)
	return err
}

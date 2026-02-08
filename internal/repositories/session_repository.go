package repositories

import (
	"cinema-system/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type SessionRepository struct {
	collection *mongo.Collection
}

func NewSessionRepository(db *mongo.Database) *SessionRepository {
	return &SessionRepository{
		collection: db.Collection("sessions"),
	}
}

func (r *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	result, err := r.collection.InsertOne(ctx, session)
	if err != nil {
		return err
	}
	session.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *SessionRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Session, error) {
	var session models.Session
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) GetByMovie(ctx context.Context, movieID primitive.ObjectID) ([]models.Session, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"movie_id":   movieID,
		"start_time": bson.M{"$gte": time.Now()},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []models.Session
	if err = cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionRepository) GetUpcoming(ctx context.Context) ([]models.Session, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"start_time": bson.M{"$gte": time.Now()},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []models.Session
	if err = cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

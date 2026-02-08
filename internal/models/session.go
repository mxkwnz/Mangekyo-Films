package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Session struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	MovieID   primitive.ObjectID `json:"movie_id" bson:"movie_id"`
	HallID    primitive.ObjectID `json:"hall_id" bson:"hall_id"`
	StartTime time.Time          `json:"start_time" bson:"start_time"`
	EndTime   time.Time          `json:"end_time" bson:"end_time"`
	Price     float64            `json:"price" bson:"price"`
}

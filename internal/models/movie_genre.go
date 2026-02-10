package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type MovieGenre struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	MovieID primitive.ObjectID `json:"movie_id" bson:"movie_id"`
	GenreID primitive.ObjectID `json:"genre_id" bson:"genre_id"`
}

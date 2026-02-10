package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name        string               `json:"name" bson:"name"`
	AgeRating   string               `json:"age_rating" bson:"age_rating"`
	Duration    int                  `json:"duration" bson:"duration"`
	Description string               `json:"description" bson:"description"`
	PosterURL   string               `json:"poster_url" bson:"poster_url"`
	Rating      float64              `json:"rating" bson:"rating"`
	Genres      []primitive.ObjectID `json:"genres" bson:"genres"`
	CreatedAt   time.Time            `json:"created_at" bson:"created_at"`
}

type Genre struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
}

type Review struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	MovieID    primitive.ObjectID `json:"movie_id" bson:"movie_id"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	Rating     int                `json:"rating" bson:"rating"`
	Comment    string             `json:"comment" bson:"comment"`
	UserName   string             `json:"user_name" bson:"user_name"`
	MovieTitle string             `json:"movie_title" bson:"movie_title"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}

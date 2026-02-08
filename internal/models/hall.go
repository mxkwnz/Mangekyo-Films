package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hall struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Location    string             `json:"location" bson:"location"`
	TotalRows   int                `json:"total_rows" bson:"total_rows"`
	SeatsPerRow int                `json:"seats_per_row" bson:"seats_per_row"`
}

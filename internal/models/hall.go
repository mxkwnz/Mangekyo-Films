package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type HallType string

const (
	HallType3D     HallType = "3D"
	HallTypeVIP    HallType = "VIP"
	HallTypeIMAX   HallType = "IMAX"
	HallTypeCommon HallType = "COMMON"
)

type Hall struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Type        HallType           `json:"type" bson:"type"`
	Location    string             `json:"location" bson:"location"`
	TotalRows   int                `json:"total_rows" bson:"total_rows"`
	SeatsPerRow int                `json:"seats_per_row" bson:"seats_per_row"`
}

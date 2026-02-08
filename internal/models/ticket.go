package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TicketStatus string

const (
	TicketBooked    TicketStatus = "BOOKED"
	TicketPaid      TicketStatus = "PAID"
	TicketCancelled TicketStatus = "CANCELLED"
)

type Ticket struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SessionID  primitive.ObjectID `json:"session_id" bson:"session_id"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	RowNumber  int                `json:"row_number" bson:"row_number"`
	SeatNumber int                `json:"seat_number" bson:"seat_number"`
	Status     TicketStatus       `json:"status" bson:"status"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}

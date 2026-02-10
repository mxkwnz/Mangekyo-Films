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

type TicketType string

const (
	TicketAdult   TicketType = "ADULT"
	TicketStudent TicketType = "STUDENT"
	TicketKid     TicketType = "KID"
	TicketPension TicketType = "PENSION"
)

type Ticket struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	SessionID  primitive.ObjectID `json:"session_id" bson:"session_id"`
	PaymentID  primitive.ObjectID `json:"payment_id" bson:"payment_id"`
	RowNumber  int                `json:"row_number" bson:"row_number"`
	SeatNumber int                `json:"seat_number" bson:"seat_number"`
	Type       TicketType         `json:"type" bson:"type"`
	Price      float64            `json:"price" bson:"price"`
	MovieTitle string             `json:"movie_title" bson:"movie_title"`
	Status     TicketStatus       `json:"status" bson:"status"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}

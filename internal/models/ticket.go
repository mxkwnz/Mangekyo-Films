package models

import "time"

type Ticket struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	SessionID  int       `json:"session_id"`
	RowNumber  int       `json:"row_number"`
	SeatNumber int       `json:"seat_number"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

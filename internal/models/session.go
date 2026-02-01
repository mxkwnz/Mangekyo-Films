package models

import "time"

type Session struct {
	ID        int       `json:"id"`
	MovieID   int       `json:"movie_id"`
	HallID    int       `json:"hall_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Price     int       `json:"price"`
}

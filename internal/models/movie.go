package models

import "time"

type Movie struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Duration    int       `json:"duration"`
	Description string    `json:"description"`
	PosterURL   string    `json:"poster_url"`
	Authors     []string  `json:"authors"`
	Rating      float64   `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
}

package models

type Hall struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	TotalRows   int    `json:"total_rows"`
	SeatsPerRow int    `json:"seats_per_row"`
}

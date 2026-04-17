package data

import "time"

type List struct {
	Name      string    `json:"name"`
	Items     []Item    `json:"items"`
	CreatedAt time.Time `json:"created_at"`
}

type Item struct {
	ID           string         `json:"id"`
	Position     int            `json:"position"` // 0 if empty
	Title        string         `json:"title"`
	By           []string       `json:"by"`
	AvgRating    float64        `json:"rating_avg"`
	TotalRatings int64          `json:"total_ratings"`
	Images       []string       `json:"images"`
	Metadata     map[string]any `json:"metadata"`
}

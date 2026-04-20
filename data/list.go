package data

import "time"

type List struct {
	ID          string    `json:"id"`
	Source      string    `json:"source"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Items       []Item    `json:"items"`
	CreatedAt   time.Time `json:"created_at"`
}

func (list List) ListType() string {
	switch list.Source {
	case "goodreads":
		return "Books"
	}

	return "Unknown"
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

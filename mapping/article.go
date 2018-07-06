package mapping

import (
	"time"
)

// Article serialization structure for indexing
type Article struct {
	SourceID    string    `json:"sourceId"`
	SourceName  string    `json:"sourceName"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	ImageURL    string    `json:"imageUrl"`
	PublishedAt time.Time `json:"publishedAt"`
	CreatedAt   time.Time `json:"createdAt"`
}

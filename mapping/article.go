package mapping

import (
	"time"
)

// Article serialization structure for indexing
type Article struct {
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	URL         string    `json:"url"`
	Source      string    `json:"source"`
	Category    string    `json:"category"`
	PublishedAt time.Time `json:"publishedAt"`
	ImportedAt  time.Time `json:"importedAt"`
}

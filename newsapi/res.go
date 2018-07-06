package newsapi

import (
	"time"
)

// Sources response from NewsAPI.org
type Sources struct {
	Status  string
	Sources []Source
}

// Source from NewsAPI.org
type Source struct {
	ID      string
	Name    string
	Country string
}

// Articles response from NewsAPI.org
type Articles struct {
	Status       string
	TotalResults int
	Articles     []Article
}

// Article from NewsAPI.org
type Article struct {
	Source
	Author      string
	Title       string
	Description string
	URL         string
	URLToImage  string
	PublishedAt time.Time
}

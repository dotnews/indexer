package newsapi_test

import (
	"testing"
	"time"

	"github.com/inthenews/indexer/config"
	"github.com/inthenews/indexer/newsapi"

	"github.com/stretchr/testify/assert"
)

var n = newsapi.New(config.New("../"))

func TestSourcesEndpoint(t *testing.T) {
	sources, err := n.Sources("en")

	assert.Nil(t, err)
	assert.Equal(t, "ok", sources.Status)
}

func TestEverythingEndpoint(t *testing.T) {
	articles, err := n.Everything(&newsapi.EverythingParams{
		Sources:  []string{"technology"},
		Date:     time.Now(),
		PageSize: 1,
		Page:     1,
	})

	assert.Nil(t, err)
	assert.Equal(t, "ok", articles.Status)
}

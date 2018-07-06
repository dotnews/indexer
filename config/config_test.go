package config_test

import (
	"os"
	"testing"

	"github.com/inthenews/indexer/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	c := config.New("../")

	assert.Equal(t, "test", c.Env)
	assert.Equal(t, "../", c.Root)
	assert.Equal(t, "http://localhost:9200", c.Elastic.URL)
	assert.Equal(t, "elastic", c.Elastic.Username)
	assert.Equal(t, "changeme", c.Elastic.Password)
	assert.NotEmpty(t, c.NewsAPI.Key)
	assert.Equal(t, "article_test", c.Article.Index)
}

func TestOverride(t *testing.T) {
	orig := os.Getenv("CONFIG")
	os.Setenv("CONFIG", "config/config.sample.json")
	c := config.New("../")

	assert.Equal(t, "sample", c.Env)
	assert.Equal(t, "../", c.Root)
	assert.Equal(t, "http://localhost:9200", c.Elastic.URL)
	assert.Equal(t, "elastic", c.Elastic.Username)
	assert.Equal(t, "changeme", c.Elastic.Password)
	assert.Empty(t, c.NewsAPI.Key)
	assert.Equal(t, "article", c.Article.Index)

	os.Setenv("CONFIG", orig)
}

package config_test

import (
	"os"
	"testing"

	"github.com/dotnews/indexer/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	c := config.New("../")

	assert.Equal(t, "test", c.Env)
	assert.Equal(t, "../", c.Root)
	assert.Equal(t, "http://localhost:9200", c.Elastic.URL)
	assert.Equal(t, "elastic", c.Elastic.Username)
	assert.Equal(t, "changeme", c.Elastic.Password)
	assert.Equal(t, "article_test", c.Article.Index)
	assert.Equal(t, "article", c.Article.Type)
	assert.Equal(t, "article.json", c.Article.Mapping)
}

func TestOverride(t *testing.T) {
	orig := os.Getenv("CONFIG")
	os.Setenv("CONFIG", "config/dev.config.json")
	c := config.New("../")

	assert.Equal(t, "dev", c.Env)
	assert.Equal(t, "../", c.Root)
	assert.Equal(t, "http://localhost:9200", c.Elastic.URL)
	assert.Equal(t, "elastic", c.Elastic.Username)
	assert.Equal(t, "changeme", c.Elastic.Password)
	assert.Equal(t, "article_v1", c.Article.Index)
	assert.Equal(t, "article", c.Article.Type)
	assert.Equal(t, "article.json", c.Article.Mapping)

	os.Setenv("CONFIG", orig)
}

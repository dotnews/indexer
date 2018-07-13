package worker

import (
	"github.com/dotnews/indexer/config"
	"github.com/dotnews/indexer/elastic"
	"github.com/dotnews/indexer/mapping"
	"github.com/golang/glog"
)

// Worker struct
type Worker struct {
	Elastic *elastic.Elastic
	Config  *config.Config
}

// New creates a new worker
func New(c *config.Config) *Worker {
	return &Worker{
		Elastic: elastic.New(c),
		Config:  c,
	}
}

// Index news articles into Elasticsearch
func (w *Worker) Index() error {
	err := w.ensureIndex()

	if err != nil {
		glog.Errorf("Failed ensuring index exists: %v", err)
		return err
	}

	err = w.Elastic.Index(
		w.Config.Article.Index,
		w.Config.Article.Type,
		"id",
		w.parseArticle(),
	)

	if err != nil {
		glog.Errorf("Failed indexing article: %v", err)
		return err
	}

	return nil
}

// Delete article index from Elastic
func (w *Worker) Delete() error {
	err := w.Elastic.Delete(w.Config.Article.Index)
	if err != nil {
		glog.Errorf("Failed deleting index: %v", err)
		return err
	}
	return nil
}

func (w *Worker) ensureIndex() error {
	mapping, err := w.Elastic.GetMapping(w.Config.Article.Mapping)
	if err != nil {
		glog.Errorf("Failed loading index mapping: %v", err)
		return err
	}

	err = w.Elastic.Ensure(w.Config.Article.Index, mapping)
	if err != nil {
		glog.Errorf("Failed ensuring index exists: %v", err)
		return err
	}

	return nil
}

func (w *Worker) parseArticle() *mapping.Article {
	return &mapping.Article{}
}

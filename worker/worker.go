package worker

import (
	"math"
	"time"

	"github.com/golang/glog"
	"github.com/inthenews/indexer/config"
	"github.com/inthenews/indexer/elastic"
	"github.com/inthenews/indexer/mapping"
	"github.com/inthenews/indexer/newsapi"
)

const (
	language        = "en"
	articleType     = "article"
	sourceSliceSize = 20
	articlePageSize = 100
)

// Worker struct
type Worker struct {
	NewsAPI *newsapi.NewsAPI
	Elastic *elastic.Elastic
	Config  *config.Config
}

// New creates a new worker
func New(c *config.Config) *Worker {
	return &Worker{
		NewsAPI: newsapi.New(c),
		Elastic: elastic.New(c),
		Config:  c,
	}
}

// Index articles from NewsAPI.org into Elastic
func (w *Worker) Index() error {
	err := w.ensureIndex()
	if err != nil {
		glog.Errorf("Failed ensuring index exists: %v", err)
		return err
	}

	sourceSlices, err := w.getSourceSlices()
	if err != nil {
		glog.Errorf("Failed getting source slices: %v", err)
		return err
	}

	err = w.indexArticles(sourceSlices, time.Now(), 0, 1)
	if err != nil {
		glog.Errorf("Failed indexing articles: %v", err)
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
	mapping, err := w.Elastic.GetMapping("article.json")
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

func (w *Worker) getSourceSlices() ([][]newsapi.Source, error) {
	s, err := w.NewsAPI.Sources(language)
	if err != nil {
		glog.Errorf("Failed fetching sources: %v", err)
		return nil, err
	}

	slices := [][]newsapi.Source{}
	for i := 0; i < len(s.Sources); i += sourceSliceSize {
		j := int(math.Min(float64(i+sourceSliceSize), float64(len(s.Sources))))
		slices = append(slices, s.Sources[i:j])
	}

	return slices, nil
}

func (w *Worker) indexArticles(sourceSlices [][]newsapi.Source, date time.Time, count, page int) error {
	sourceIds := []string{}
	for _, source := range sourceSlices[0] {
		sourceIds = append(sourceIds, source.ID)
	}

	a, err := w.NewsAPI.Everything(&newsapi.EverythingParams{
		Sources:  sourceIds,
		Date:     date,
		PageSize: articlePageSize,
		Page:     page,
	})

	if err != nil {
		glog.Errorf("Failed fetching articles: %v", err)
		return err
	}

	for _, article := range a.Articles {
		err = w.Elastic.Index(
			w.Config.Article.Index,
			articleType,
			w.parseArticle(article, date),
		)

		if err != nil {
			glog.Errorf("Failed indexing article: %v", err)
			return err
		}
	}

	count += len(a.Articles)
	glog.Infof("Indexed %d/%d articles. Source slice: %v", count, a.TotalResults, sourceSlices[0])

	if count < a.TotalResults {
		err := w.indexArticles(sourceSlices, date, count, page+1)
		if err != nil {
			glog.Errorf("Failed recursing into next page: %v", err)
			return err
		}
	}

	if len(sourceSlices) > 1 {
		err := w.indexArticles(sourceSlices[1:], date, 0, 1)
		if err != nil {
			glog.Errorf("Failed recursing into next source slice: %v", err)
			return err
		}
	}

	return nil
}

func (w *Worker) parseArticle(article newsapi.Article, date time.Time) mapping.Article {
	return mapping.Article{
		SourceID:    article.Source.ID,
		SourceName:  article.Source.Name,
		Author:      article.Author,
		Title:       article.Title,
		Description: article.Description,
		URL:         article.URL,
		ImageURL:    article.URLToImage,
		PublishedAt: article.PublishedAt,
		CreatedAt:   date,
	}
}

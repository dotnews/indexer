package elastic

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/golang/glog"
	"github.com/inthenews/indexer/config"
	"github.com/olivere/elastic"
)

const (
	// ElasticErrorNotFound from get query
	ElasticErrorNotFound = "Error 404 (Not Found)"

	elasticScheme   = "https"
	elasticSniffing = false
)

// Elastic client
type Elastic struct {
	Context context.Context
	Client  *elastic.Client
	Config  *config.Config
}

// New create a new elastic client
func New(c *config.Config) *Elastic {
	cli, err := elastic.NewClient(
		elastic.SetURL(c.Elastic.URL),
		elastic.SetBasicAuth(c.Elastic.Username, c.Elastic.Password),
		elastic.SetScheme(elasticScheme),
		elastic.SetSniff(elasticSniffing),
	)

	if err != nil {
		glog.Fatalf("Failed creating indexer: %v", err)
	}

	return &Elastic{
		Context: context.Background(),
		Client:  cli,
		Config:  c,
	}
}

// Create index
func (el *Elastic) Create(index, mapping string) error {
	res, err := el.Client.CreateIndex(index).BodyString(mapping).Do(el.Context)

	if err != nil {
		glog.Errorf("Failed creating index: %v", err)
		return err
	}

	if !res.Acknowledged {
		return fmt.Errorf("Index creation was not acknowledged: %v", res)
	}

	return nil
}

// Ensure index exists
func (el *Elastic) Ensure(index, mapping string) error {
	exists, err := el.Client.IndexExists(index).Do(el.Context)

	if err != nil {
		glog.Errorf("Failed checking if index exists: %v", err)
		return err
	}

	if !exists {
		return el.Create(index, mapping)
	}

	return nil
}

// Delete index
func (el *Elastic) Delete(index string) error {
	res, err := el.Client.DeleteIndex(index).Do(el.Context)

	if err != nil {
		glog.Errorf("Failed deleting index: %v", err)
		return err
	}

	if !res.Acknowledged {
		return fmt.Errorf("Index deletion was not acknowledged: %v", res)
	}

	return nil
}

// Index a document
func (el *Elastic) Index(name, typ string, body interface{}) error {
	_, err := el.Client.Index().
		Index(name).
		Type(typ).
		BodyJson(body).
		Do(el.Context)

	if err != nil {
		glog.Errorf("Failed indexing document: %v", err)
		return err
	}

	return nil
}

// Get document source
func (el *Elastic) Get(index, typ, id string) ([]byte, error) {
	res, err := el.Client.Get().
		Index(index).
		Type(typ).
		Id(id).
		Do(el.Context)

	if err != nil {
		if strings.Contains(err.Error(), ElasticErrorNotFound) {
			return nil, err
		}
		glog.Errorf("Failed getting document: %v", err)
		return nil, err
	}

	return *res.Source, err
}

// GetMapping for index
func (el *Elastic) GetMapping(filename string) (string, error) {
	f := el.Config.GetPath(fmt.Sprintf("mapping/%s", filename))
	b, err := ioutil.ReadFile(f)

	if err != nil {
		glog.Errorf("Failed reading '%s' with error: %v", f, err)
		return "", err
	}

	return string(b), nil
}

// Flush indices
func (el *Elastic) Flush(indices ...string) error {
	_, err := el.Client.Flush(indices...).Do(el.Context)

	if err != nil {
		glog.Errorf("Failed flushing indices '%v': %v", indices, err)
		return err
	}

	return nil
}

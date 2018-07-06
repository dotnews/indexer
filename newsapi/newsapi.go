package newsapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/inthenews/indexer/config"
)

const endpoint = "https://newsapi.org/v2"

// NewsAPI client
type NewsAPI struct {
	Client *http.Client
	Config *config.Config
}

// EverythingParams for NewsAPI endpoint
type EverythingParams struct {
	Sources  []string
	Date     time.Time
	PageSize int
	Page     int
}

// New creates a new client
func New(c *config.Config) *NewsAPI {
	return &NewsAPI{
		Client: &http.Client{},
		Config: c,
	}
}

// Fetch url with query params from NewsAPI.org
func (n *NewsAPI) Fetch(url string, query *url.Values) ([]byte, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		endpoint+url,
		nil,
	)

	if err != nil {
		glog.Errorf("Failed creating request: %v", err)
		return nil, err
	}

	query.Add("apiKey", n.Config.NewsAPI.Key)
	req.URL.RawQuery = query.Encode()

	res, err := n.Client.Do(req)
	if err != nil {
		glog.Errorf("Failed processing request: %v", err)
		return nil, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		glog.Errorf("Failed reading response body: %v", err)
		return nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf(
			"Request failed: %d: %s: %s",
			res.StatusCode,
			res.Status,
			string(resBody),
		)
	}

	return resBody, nil
}

// Sources endpoint from NewsAPI.org
func (n *NewsAPI) Sources(language string) (*Sources, error) {
	query := &url.Values{}
	query.Add("language", language)

	b, err := n.Fetch("/sources", query)
	if err != nil {
		glog.Errorf("Failed fetching sources: %v Query: %+v", err, query)
		return nil, err
	}

	var sources Sources
	err = json.Unmarshal(b, &sources)
	if err != nil {
		glog.Errorf("Failed parsing sources: %v", err)
		return nil, err
	}

	return &sources, nil
}

// Everything endpoint from NewsAPI.org
func (n *NewsAPI) Everything(params *EverythingParams) (*Articles, error) {
	query := &url.Values{}
	date := params.Date.Format(time.RFC3339)[:10]
	query.Add("sources", strings.Join(params.Sources, ","))
	query.Add("from", date)
	query.Add("to", date)
	query.Add("pageSize", strconv.Itoa(params.PageSize))
	query.Add("page", strconv.Itoa(params.Page))

	b, err := n.Fetch("/everything", query)
	if err != nil {
		glog.Errorf("Failed fetching articles from everything endpoint: %v Query: %+v", err, query)
		return nil, err
	}

	var articles Articles
	err = json.Unmarshal(b, &articles)
	if err != nil {
		glog.Errorf("Failed parsing articles from everything endpoint: %v", err)
		return nil, err
	}

	return &articles, nil
}

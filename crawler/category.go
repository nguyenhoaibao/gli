package crawler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

type CategoryRenderer interface {
	Total() int
	ItemN(i int) string
	Render() io.Reader
}

type CategoryParser interface {
	Parse(io.Reader, int) (CategoryRenderer, error)
}

type categoryCrawler struct {
	name            string
	url             string
	limit           int
	cachedInSeconds time.Duration
	mu              sync.Mutex
	cached          CategoryRenderer
}

var categoryParsers = make(map[string]CategoryParser)

func NewCategoryCrawler(name string, url string, limit int, cachedInSeconds time.Duration) *categoryCrawler {
	c := &categoryCrawler{
		name:            name,
		url:             url,
		limit:           limit,
		cachedInSeconds: cachedInSeconds,
	}
	c.refreshCached()

	return c
}

func RegisterCategoryParser(name string, p CategoryParser) error {
	if _, exists := categoryParsers[name]; exists {
		return fmt.Errorf("Parser %s was already registered", name)
	}
	categoryParsers[name] = p
	return nil
}

func (c *categoryCrawler) Crawl() (CategoryRenderer, error) {
	if result := c.getCached(); result != nil {
		return result, nil
	}

	resp, err := c.Download()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, err := c.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	c.cache(result)
	return result, nil
}

func mockServer(content string) *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")
		w.WriteHeader(200)
		fmt.Fprintf(w, content)
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

func (c *categoryCrawler) Download() (*http.Response, error) {
	return http.Get(c.url)
	// content, err := ioutil.ReadFile(filepath.Join("testdata", c.name+".html"))
	// if err != nil {
	// 	return nil, err
	// }
	//
	// server := mockServer(string(content[:]))
	// defer server.Close()
	//
	// return http.Get(server.URL)
}

func (c *categoryCrawler) Parse(r io.Reader) (CategoryRenderer, error) {
	p, exists := categoryParsers[c.name]
	if !exists {
		return nil, fmt.Errorf("Parser %s does not exist", c.name)
	}
	return p.Parse(r, c.limit)
}

func (c *categoryCrawler) getCached() CategoryRenderer {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cached
}

func (c *categoryCrawler) cache(result CategoryRenderer) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cached = result
}

func (c *categoryCrawler) refreshCached() {
	ticker := time.NewTicker(time.Second * c.cachedInSeconds)
	go func() {
		for {
			<-ticker.C
			c.mu.Lock()
			c.cached = nil
			c.mu.Unlock()
		}
	}()
}

package crawler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

type ItemsRenderer interface {
	Total() int
	ItemN(i int) string
	Render() io.Reader
}

type ItemsParser interface {
	Parse(io.Reader, int) (ItemsRenderer, error)
}

type itemsCrawler struct {
	name            string
	url             string
	limit           int
	cachedInSeconds time.Duration
	mu              sync.Mutex
	cached          ItemsRenderer
}

var itemsParser = make(map[string]ItemsParser)

func NewItemsCrawler(name string, url string, limit int, cachedInSeconds time.Duration) *itemsCrawler {
	ic := &itemsCrawler{
		name:            name,
		url:             url,
		limit:           limit,
		cachedInSeconds: cachedInSeconds,
	}
	ic.refreshCached()

	return ic
}

func RegisterItemsParser(name string, p ItemsParser) error {
	if _, exists := itemsParser[name]; exists {
		return fmt.Errorf("Parser %s was already registered", name)
	}
	itemsParser[name] = p
	return nil
}

func (c *itemsCrawler) Crawl() (ItemsRenderer, error) {
	if items := c.getCached(); items != nil {
		return items, nil
	}

	resp, err := c.Download()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	items, err := c.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	c.cache(items)
	return items, nil
}

func mockServer(content string) *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")
		w.WriteHeader(200)
		fmt.Fprintf(w, content)
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

func (c *itemsCrawler) Download() (*http.Response, error) {
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

func (c *itemsCrawler) Parse(r io.Reader) (ItemsRenderer, error) {
	p, exists := itemsParser[c.name]
	if !exists {
		return nil, fmt.Errorf("Parser %s does not exist", c.name)
	}
	return p.Parse(r, c.limit)
}

func (c *itemsCrawler) getCached() ItemsRenderer {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cached
}

func (c *itemsCrawler) cache(items ItemsRenderer) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cached = items
}

func (c *itemsCrawler) refreshCached() {
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

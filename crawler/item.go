package crawler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

type ItemRenderer interface {
	Render() io.Reader
}

type ItemParser interface {
	Parse(io.Reader) (ItemRenderer, error)
}

type itemCrawler struct {
	name       string
	urlPattern string
	cached     map[string]ItemRenderer
}

var itemParsers = make(map[string]ItemParser)

func NewItemCrawler(name string, urlPattern string) *itemCrawler {
	if urlPattern == "" {
		return nil
	}
	return &itemCrawler{
		name:       name,
		urlPattern: urlPattern,
		cached:     make(map[string]ItemRenderer),
	}
}

func mockServer(content string) *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		// w.Header().Set("Content-type", "text/html")
		fmt.Fprintf(w, content)
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

func RegisterItemParser(name string, p ItemParser) error {
	if _, exists := itemParsers[name]; exists {
		return fmt.Errorf("Parser %s was already registered", name)
	}
	itemParsers[name] = p
	return nil
}

func (c *itemCrawler) Crawl(id string) (ItemRenderer, error) {
	if item := c.getCached(id); item != nil {
		return item, nil
	}

	url := fmt.Sprintf(c.urlPattern, id)
	resp, err := c.Download(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	item, err := c.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	c.cache(id, item)
	return item, nil
}

func (c *itemCrawler) Download(url string) (*http.Response, error) {
	// return http.Get(url)
	// content, err := ioutil.ReadFile("./testdata/markdown_content.md")
	// if err != nil {
	// 	return nil, err
	// }

	server := mockServer("")
	defer server.Close()

	return http.Get(server.URL)
}

func (c *itemCrawler) Parse(r io.Reader) (ItemRenderer, error) {
	p, exists := itemParsers[c.name]
	if !exists {
		return nil, fmt.Errorf("Parser %s does not exist", c.name)
	}
	return p.Parse(r)
}

func (c *itemCrawler) getCached(id string) ItemRenderer {
	item, exists := c.cached[id]
	if !exists {
		return nil
	}
	return item
}

func (c *itemCrawler) cache(id string, item ItemRenderer) {
	c.cached[id] = item
}

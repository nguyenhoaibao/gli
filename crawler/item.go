package crawler

import (
	"fmt"
	"io"
	"net/http"
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

var itemParser = make(map[string]ItemParser)

func NewItemCrawler(name string, urlPattern string) *itemCrawler {
	return &itemCrawler{
		name:       name,
		urlPattern: urlPattern,
		cached:     make(map[string]ItemRenderer),
	}
}

func RegisterItemParser(name string, p ItemParser) error {
	if _, exists := itemParser[name]; exists {
		return fmt.Errorf("Parser %s was already registered", name)
	}
	itemParser[name] = p
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
	return http.Get(url)
	// content, err := ioutil.ReadFile(filepath.Join("testdata", c.name+".json"))
	// if err != nil {
	// 	return nil, err
	// }
	//
	// server := mockServer(string(content[:]))
	// defer server.Close()
	//
	// return http.Get(server.URL)
}

func (c *itemCrawler) Parse(r io.Reader) (ItemRenderer, error) {
	p, exists := itemParser[c.name]
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
package crawler

import (
	"fmt"

	"github.com/nguyenhoaibao/gli/app"
)

type jsonCrawler struct {
	site *app.Site
}

func NewJsonCrawler(site *app.Site) *htmlCrawler {
	return &htmlCrawler{site: site}
}

func (c *jsonCrawler) Crawl() (Renderer, error) {
	name := c.site.Name
	url := c.site.Url

	resp, err := Query(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	parser, exists := parsers[name]
	if !exists {
		return nil, fmt.Errorf("Parser %s does not exist", name)
	}

	return parser.Parse(resp.Body, c.site.Limit)
}

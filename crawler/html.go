package crawler

import (
	"errors"
	"io"

	"github.com/nguyenhoaibao/gli/app"
)

type htmlCrawler struct {
	site *app.Site
}

func NewHtmlCrawler(site *app.Site) *htmlCrawler {
	return &htmlCrawler{site: site}
}

func (c *htmlCrawler) Crawl() (io.Writer, error) {
	name := c.site.Name
	url := c.site.Url

	resp, err := Query(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	parser, exists := parsers[name]
	if !exists {
		return nil, errors.New("Parser does not exist")
	}

	return parser.Parse(resp.Body)
}

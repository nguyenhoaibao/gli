package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"

	"github.com/nguyenhoaibao/gli/app"
)

type htmlCrawler struct {
	site *app.Site
}

func NewHtmlCrawler(site *app.Site) *htmlCrawler {
	return &htmlCrawler{site: site}
}

func mockServer(content string) *httptest.Server {
	handerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")
		w.WriteHeader(200)
		fmt.Fprintf(w, content)
	}

	return httptest.NewServer(http.HandlerFunc(handerFunc))
}

func (c *htmlCrawler) Crawl() (Renderer, error) {
	name := c.site.Name
	// url := c.site.Url

	content, err := ioutil.ReadFile(filepath.Join("parsers/testdata", name+".html"))
	server := mockServer(string(content[:]))
	defer server.Close()

	resp, err := Query(server.URL)
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

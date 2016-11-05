package crawler_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/nguyenhoaibao/gli/app"
	"github.com/nguyenhoaibao/gli/crawler"
	_ "github.com/nguyenhoaibao/gli/parsers"
)

func mockServer(content string) *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")
		w.WriteHeader(200)
		fmt.Fprintf(w, content)
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

func TestCrawl(t *testing.T) {
	var sites = []*app.Site{
		&app.Site{
			Name:  "github_trending",
			Limit: 10,
		},
	}

	for _, site := range sites {
		content, err := ioutil.ReadFile(filepath.Join("../testdata", site.Name+".html"))
		if err != nil {
			t.Error(err)
		}

		server := mockServer(string(content[:]))
		defer server.Close()

		c := crawler.New(site.Name, server.URL, site.Limit)
		items, err := c.Crawl()
		if err != nil {
			t.Error(err)
		}
		fmt.Println(items.Render())
	}
}

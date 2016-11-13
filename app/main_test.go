package app_test

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
			Name: "github",
			Types: []*app.Items{
				&app.Items{
					Name:            "trending",
					Limit:           25,
					CachedInSeconds: 300,
				},
			},
		},
	}

	for _, site := range sites {
		for _, items := range site.Types {
			name := fmt.Sprintf("%s_%s", site.Name, items.Name)
			itemsContent, err := ioutil.ReadFile(filepath.Join("../testdata", fmt.Sprintf("%s.html", name)))
			if err != nil {
				t.Error(err)
			}

			server := mockServer(string(itemsContent[:]))
			defer server.Close()

			c := crawler.NewItemsCrawler(name, server.URL, items.Limit, items.CachedInSeconds)
			results, err := c.Crawl()
			if err != nil {
				t.Error(err)
			}
			if results == nil {
				t.Fatal("Items should not be nil")
			}
			if results.Total() != items.Limit {
				t.Fatal("Total items do not match limit")
			}
		}
	}
}

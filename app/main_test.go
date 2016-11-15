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
			Categories: []*app.Category{
				&app.Category{
					Name:            "trending",
					Limit:           25,
					CachedInSeconds: 300,
				},
			},
		},
	}

	for _, s := range sites {
		for _, c := range s.Categories {
			cName := fmt.Sprintf("%s_%s", s.Name, c.Name)
			cContent, err := ioutil.ReadFile(filepath.Join("../testdata", fmt.Sprintf("%s.html", cName)))
			if err != nil {
				t.Error(err)
			}

			server := mockServer(string(cContent[:]))
			defer server.Close()

			cCrawler := crawler.NewCategoryCrawler(cName, server.URL, c.Limit, c.CachedInSeconds)
			results, err := cCrawler.Crawl()
			if err != nil {
				t.Error(err)
			}
			if results == nil {
				t.Fatal("Items should not be nil")
			}
			if results.Total() != c.Limit {
				t.Fatal("Total items do not match limit")
			}
		}
	}
}

package parsers_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/nguyenhoaibao/gli/app"
	"github.com/nguyenhoaibao/gli/crawler"
	"github.com/nguyenhoaibao/gli/parsers"
)

func mockServer(content string) *httptest.Server {
	handerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")
		w.WriteHeader(200)
		fmt.Fprintf(w, content)
	}

	return httptest.NewServer(http.HandlerFunc(handerFunc))
}

func TestGithubTrendingParseResults(t *testing.T) {
	var site = &app.Site{
		Type: "html",
		Name: "github_trending",
	}

	content, err := ioutil.ReadFile(filepath.Join("testdata", site.Name+".html"))
	if err != nil {
		t.Error(err)
	}

	server := mockServer(string(content[:]))
	defer server.Close()

	site.Url = server.URL

	resp, err := crawler.Query(site.Url)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	doc, err := crawler.GetDocumentFromReader(resp.Body)

	p := parsers.NewGHTrendingParser()
	items := p.ParseItems(doc, 10)
	if len(items) <= 0 {
		t.Errorf("Cannot get any results from site %s", site.Name)
	}

	items.Render()
}

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

var site = &app.Site{
	Type: "html",
	Name: "github_trending",
}

func mockServer(content string) *httptest.Server {
	handerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")
		w.WriteHeader(200)
		fmt.Fprintf(w, content)
	}

	return httptest.NewServer(http.HandlerFunc(handerFunc))
}

func TestGetResults(t *testing.T) {
	content, err := ioutil.ReadFile(filepath.Join("testdata", site.Name+".html"))
	if err != nil {
		t.Fatal(err)
	}

	server := mockServer(string(content[:]))
	defer server.Close()

	site.Url = server.URL

	resp, err := crawler.Query(site.Url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := crawler.GetDocumentFromReader(resp.Body)

	var parser parsers.GithubTrendingParser
	results := parser.ParseResults(doc)
	if len(results) == 0 {
		t.Errorf("Cannot get any results from site %s", site.Name)
	}

	parser.Display(results)
}

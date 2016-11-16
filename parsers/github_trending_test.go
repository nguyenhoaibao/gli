package parsers

import (
	"bufio"
	"os"
	"testing"

	"github.com/nguyenhoaibao/gli/crawler"
)

func TestGithubTrendingParser(t *testing.T) {
	f, err := os.Open("../testdata/github_trending.html")
	if err != nil {
		t.Error(err)
	}

	var (
		p     ghTrendingParser
		r     = bufio.NewReader(f)
		limit = 10
	)

	doc, err := crawler.GetDocumentFromReader(r)
	if err != nil {
		t.Error(err)
	}

	results := p.parse(doc, limit)
	if len(results) == 0 {
		t.Error("Cannot parse any items")
	}
	for _, r := range results {
		if r.Title == "" {
			t.Error("Cannot parse title")
		}
		if r.Owner == "" {
			t.Error("Cannot parse owner")
		}
		if r.Description == "" {
			t.Error("Cannot parse description")
		}
		if r.TodayStars == "" {
			t.Error("Cannot parse today stars")
		}
	}
}

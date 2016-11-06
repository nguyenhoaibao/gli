package parsers

import (
	"bufio"
	"os"
	"testing"

	"github.com/nguyenhoaibao/gli/crawler"
)

func TestGithubTrendingParseItems(t *testing.T) {
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

	items := p.parse(doc, limit)
	if len(*items) == 0 {
		t.Error("Cannot parse any items")
	}
	for _, item := range *items {
		if item.Title == "" {
			t.Error("Cannot parse title")
		}
		if item.Owner == "" {
			t.Error("Cannot parse owner")
		}
		if item.Description == "" {
			t.Error("Cannot parse description")
		}
		if item.TodayStars == "" {
			t.Error("Cannot parse today stars")
		}
	}
}

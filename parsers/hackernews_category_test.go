package parsers

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/nguyenhoaibao/gli/crawler"
)

func TestHackernewsCategoryParser(t *testing.T) {
	f, err := os.Open("../testdata/hackernews_category.html")
	if err != nil {
		t.Error(err)
	}

	var (
		r     = bufio.NewReader(f)
		limit = 10
	)

	doc, err := crawler.DocumentFromReader(r)
	if err != nil {
		t.Error(err)
	}

	items := parseHNCategory(doc, limit)
	if len(items) == 0 {
		t.Error("Cannot parse any items")
	}
	fmt.Println(items.Render())
	// for _, r := range results {
	// 	if r.Title == "" {
	// 		t.Error("Cannot parse title")
	// 	}
	// 	if r.Owner == "" {
	// 		t.Error("Cannot parse owner")
	// 	}
	// 	if r.Description == "" {
	// 		t.Error("Cannot parse description")
	// 	}
	// 	if r.TodayStars == "" {
	// 		t.Error("Cannot parse today stars")
	// 	}
	// }
}
